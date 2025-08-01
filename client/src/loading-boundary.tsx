import { Fragment, useLayoutEffect, useState, useRef } from "react";
import Loading from "@/pages/loading";
import { useLocation, Outlet } from "react-router-dom";

// TODO: Make it work like Suspense fallback on every navigation instead of initial load only.
export function LoadingBoundary() {
  const { pathname } = useLocation();
  const [loaded, setLoaded] = useState(false);
  const timeoutRef = useRef<number | null>(null);

  useLayoutEffect(() => {
    const handlePageLoaded = (e: Event) => {
      const customEvent = e as CustomEvent;
      const eventPathname = customEvent.detail.pathname;

      if (eventPathname === pathname) {
        timeoutRef.current = window.setTimeout(() => {
          setLoaded(true);
        }, 100); // slight delay to allow fade
      }
    };

    window.addEventListener("PageLoaded", handlePageLoaded);

    return () => {
      window.removeEventListener("PageLoaded", handlePageLoaded);
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, [pathname]);

  return (
    <Fragment>
      {/* Overlay loading screen */}
      {!loaded && (
        <div className="fixed inset-0 overflow-auto transition-opacity duration-500 opacity-100">
          <Loading />
        </div>
      )}

      {/* Content area */}
      <section
        className={`transition-opacity duration-500 ${loaded ? "opacity-100" : "opacity-0 pointer-events-none h-screen overflow-hidden"}`}
      >
        <Outlet />
      </section>
    </Fragment>
  );
}
