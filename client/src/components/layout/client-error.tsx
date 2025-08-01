import { Fragment } from "react"
import { StaticMetadata } from "@/contexts/metadata"

type Props = {
  error: Error;
  resetErrorBoundary: (() => void) | null;
};

export default function ClientError({ error, resetErrorBoundary }: Props) {

  return (
    <Fragment>
      <StaticMetadata />

      <div className="relative flex flex-col items-center justify-center h-screen w-screen bg-background text-foreground overflow-hidden">
        {/* Content */}
        <h1 className="text-5xl md:text-6xl font-extrabold">
          Something went wrong
        </h1>
        <p className="mt-4 text-lg md:text-xl max-w-xl text-muted-foreground animate-fade-in z-10">
          This is a client-side issue. Here’s what React says:
        </p>

        <pre className="mt-4 bg-muted p-4 rounded-lg text-sm text-left max-w-xl overflow-auto z-10">
          {error.message}
        </pre>

        {resetErrorBoundary !== null && (
          <button
            onClick={resetErrorBoundary}
            className="mt-6 px-6 py-3 bg-primary text-primary-foreground rounded-lg shadow-lg hover:bg-primary/80 transition z-10"
          >
            Retry
          </button>
        )}

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
          animation: gradient 8s ease infinite;
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
      `}</style>
      </div>
    </Fragment>
  )
}
