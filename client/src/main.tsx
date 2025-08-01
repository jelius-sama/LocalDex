import { StrictMode, useState, Fragment } from 'react'
import '@/index.css'
import { BrowserRouter } from 'react-router-dom'
import { ThemeProvider } from '@/contexts/theme'
import { createRoot } from 'react-dom/client'
import { ConfigProvider } from '@/contexts/config'
import { lazy, Suspense, useLayoutEffect, useEffect, type ReactNode } from 'react'
import { Outlet, Route, Routes, useLocation } from 'react-router-dom'
import { useConfig } from "@/contexts/config"
import { QueryClientProvider, QueryClient } from '@tanstack/react-query'
import Loading from "@/components/layout/loading"
import { ErrorBoundary } from "@/error-boundary"
import { LoadingBoundary } from "@/loading-boundary"

const queryClient = new QueryClient()

const Home = lazy(() => import("@/pages/home"))
const GenericHTTPError = lazy(() => import("@/components/layout/http-error"))
const Toaster = lazy(() => import('@/components/ui/sonner'))

let rootEl = document.getElementById('root') as HTMLDivElement | null;

if (!rootEl) {
  if (process.env.NODE_ENV === "development") {
    throw new Error("Root element not found!")
  } else {
    const div = document.createElement('div');
    div.id = "root"
    document.body.appendChild(div);
    rootEl = div
  }
}

const App = () => {
  const { pathname } = useLocation();
  const priorityPaths = ["/", "/links", "/blogs", "/blog/"]

  useLayoutEffect(() => {
    document.documentElement.scrollTo({
      top: 0,
      left: 0,
      behavior: "instant",
    });
  }, [pathname]);


  const isPriorityPath = priorityPaths.some((path) => {
    if (path.endsWith("/*")) {
      // INFO: Remove the '*' but keep the '/'
      const basePath = path.slice(0, -1);
      return pathname.startsWith(basePath);
    } else {
      return pathname === path;
    }
  });

  return (
    <Fragment>
      <ErrorBoundary>
        <Suspense fallback={isPriorityPath ? null : <Loading />}>
          {isPriorityPath ? <LoadingBoundary /> : <Outlet />}
        </Suspense>
      </ErrorBoundary>
    </Fragment>
  )
}

const ServerErrorWrapper = ({ comp }: { comp: ReactNode }) => {
  const [errorPath, setErrorPath] = useState<string | null>(null)
  const { pathname } = useLocation()
  const { setSSRData } = useConfig()
  // INFO: The following state is to avoid race condition
  const [isSSRLoaded, setIsSSRLoaded] = useState(false)

  useLayoutEffect(() => {
    const script = document.getElementById('__SERVER_DATA__')
    if (script && script.textContent) {
      try {
        const data = JSON.parse(script.textContent)
        setSSRData(data)

        // Identify if it's error data (you can refine this signature check)
        if ("status" in data && data.status === 500) {
          setErrorPath(pathname)
        }

        script.remove()
      } catch (err) {
        console.error('Failed to parse SSR data:', err)
      }
    }
    setIsSSRLoaded(true)
  }, [])

  // Clear error when navigating to a different path
  useEffect(() => {
    if (errorPath && pathname !== errorPath) {
      setErrorPath(null)
    }
  }, [pathname, errorPath])

  return !isSSRLoaded ? <Loading /> : errorPath === pathname ? <GenericHTTPError error={"500 - Internal Server Error"} message={"Something broke on my end. If you’re me, fix it now. ⚡"} status={500} /> : comp
}

export const Authenticate = ({ page }: { page: React.ReactNode }) => {
  const [status, setStatus] = useState<"pending" | "success" | "error">("pending")

  useEffect(() => {
    fetch(`/api/verify_auth`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (res.status === 200) {
          setStatus("success")
        } else if (res.status === 498) {
          // INFO: Expired token
          setStatus("error")
        } else {
          setStatus("error")
        }
      })
      .catch(() => {
        setStatus("error")
      })
  }, [])

  if (status === "pending") return <Loading />
  if (status === "error") return <GenericHTTPError error={"404 - Page Not Found"} message={"This route does not exists"} status={404} />

  return page
}

const reactRoot = createRoot(rootEl);
reactRoot.render(
  <StrictMode>
    <ConfigProvider>
      <BrowserRouter>
        <QueryClientProvider client={queryClient}>
          <ThemeProvider defaultTheme="dark" storageKey="theme">
            <Routes>
              <Route path='/' element={<App />}>
                <Route path='/' element={<ServerErrorWrapper comp={<Home />} />} />
                <Route path='*' element={<ServerErrorWrapper comp={<GenericHTTPError error={"404 - Page Not Found"} message={"This route does not exists"} status={404} />} />} />
              </Route>
            </Routes>

            <Suspense><Toaster richColors={true} /></Suspense>
          </ThemeProvider>
        </QueryClientProvider>
      </BrowserRouter>
    </ConfigProvider>
  </StrictMode>
);

// Registering service worker
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/assets/sw.js')
      .then((registration) => {
        console.log('Service Worker registered with scope:', registration.scope);
      })
      .catch((err) => {
        console.error('Service Worker registration failed:', err);
      });
  });
}
