import { Fragment } from "react"
import { PathBasedMetadata } from "@/contexts/metadata"
import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom"

export default function GenericHTTPError({ error, message, status }: { error: string, message: string, status: 404 | 500 }) {
  const metadataId = status === 404 ? "#not_found" : "#internal_server_error"

  return (
    <Fragment>
      <PathBasedMetadata paths={["*", metadataId]} />

      <div className="relative flex flex-col items-center justify-center h-screen w-screen bg-background text-foreground overflow-hidden">
        {/* Animated gradient background */}
        <div className="absolute inset-0 bg-gradient-to-tr from-purple-600/20 via-pink-500/20 to-indigo-500/20 animate-gradient opacity-50"></div>

        {/* Floating SVG icons */}
        <svg
          className="absolute w-12 h-12 text-purple-400 animate-float-slow"
          style={{ top: "15%", left: "10%" }}
          fill="currentColor"
          viewBox="0 0 24 24"
        >
          <path d="M12 2L15 8H9L12 2ZM12 22L9 16H15L12 22ZM2 12L8 15V9L2 12ZM22 12L16 9V15L22 12Z" />
        </svg>
        <svg
          className="absolute w-10 h-10 text-pink-400 animate-float"
          style={{ top: "70%", left: "20%" }}
          fill="currentColor"
          viewBox="0 0 24 24"
        >
          <path d="M12 4.354a9 9 0 110 15.292A9 9 0 0112 4.354z" />
        </svg>
        <svg
          className="absolute w-14 h-14 text-indigo-400 animate-float-slow"
          style={{ top: "40%", right: "15%" }}
          fill="currentColor"
          viewBox="0 0 24 24"
        >
          <path d="M4 4h16v16H4V4z" />
        </svg>

        {/* Main Content */}
        <h1 className="text-6xl md:text-7xl font-extrabold animate-pulse z-10">{error}</h1>
        <p className="mt-4 text-xl md:text-2xl text-muted-foreground animate-fade-in z-10">
          {message}
        </p>

        <div className="flex items-center justify-center gap-x-4">
          <Button onClick={() => document.location.reload()} size="lg" className="mt-6 z-10">Retry</Button>
          {status === 404 && (
            <Button asChild={true} size="lg" className="mt-6 z-10" variant="secondary">
              <Link to="/">Return Home</Link>
            </Button>
          )}
        </div>

        {/* Animations */}
        <style>{`
        @keyframes gradient {
          0% {
            background-position: 0% 50%;
          }
          50% {
            background-position: 100% 50%;
          }
          100% {
            background-position: 0% 50%;
          }
        }
        .animate-gradient {
          background-size: 200% 200%;
          animation: gradient 10s ease infinite;
        }
        @keyframes fade-in {
          0% {
            opacity: 0;
            transform: translateY(20px);
          }
          100% {
            opacity: 1;
            transform: translateY(0);
          }
        }
        .animate-fade-in {
          animation: fade-in 1s ease forwards;
        }
        @keyframes float {
          0%,
          100% {
            transform: translateY(0) rotate(0deg);
          }
          50% {
            transform: translateY(-15px) rotate(10deg);
          }
        }
        .animate-float {
          animation: float 6s ease-in-out infinite;
        }
        .animate-float-slow {
          animation: float 10s ease-in-out infinite;
        }
      `}</style>
      </div>
    </Fragment>
  )
}
