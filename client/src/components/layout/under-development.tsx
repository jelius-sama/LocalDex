export default function UnderConstruction({ reminder }: { reminder?: string }) {
  return (
    <div className="relative flex flex-col items-center justify-center h-screen w-screen bg-background text-foreground overflow-hidden">
      {/* Animated Background Gradient */}
      <div className="absolute inset-0 bg-gradient-to-r from-purple-500/20 via-blue-500/20 to-pink-500/20 animate-gradient opacity-50"></div>

      {/* Floating shapes */}
      <div className="absolute w-72 h-72 bg-purple-500/10 rounded-full blur-3xl animate-pulse top-10 left-10"></div>
      <div className="absolute w-64 h-64 bg-pink-500/10 rounded-full blur-3xl animate-pulse bottom-10 right-10"></div>

      {/* Main Content */}
      <div className="relative z-10 flex flex-col items-center text-center p-6">
        <h1 className="text-5xl md:text-6xl font-extrabold tracking-tight animate-bounce">
          🚧 Under Construction
        </h1>
        <p className="mt-4 text-lg md:text-xl max-w-xl text-muted-foreground animate-fade-in">
          This space is still under construction—just for me. If you're not me, you're lost. 🛠️
        </p>

        {/* Animated Loader */}
        <div className="mt-8 flex space-x-2">
          <div className="w-4 h-4 bg-primary rounded-full animate-bounce delay-75"></div>
          <div className="w-4 h-4 bg-primary rounded-full animate-bounce delay-150"></div>
          <div className="w-4 h-4 bg-primary rounded-full animate-bounce delay-300"></div>
        </div>

        {reminder && reminder !== "" && (
          <div className="mt-8 mb-4 px-3 py-1 rounded-md bg-blue-600 animate-fade-in">
            <p className="">REMINDER: {reminder}</p>
          </div>
        )}
      </div>

      {/* Particle Animation using Tailwind keyframes */}
      <div className="absolute inset-0 overflow-hidden z-0">
        {Array.from({ length: 20 }).map((_, i) => (
          <span
            key={i}
            className="absolute w-1 h-1 bg-primary rounded-full animate-float"
            style={{
              top: `${Math.random() * 100}%`,
              left: `${Math.random() * 100}%`,
              animationDuration: `${5 + Math.random() * 5}s`,
              animationDelay: `${Math.random() * 5}s`,
            }}
          ></span>
        ))}
      </div>

      {/* Custom Animations */}
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
        @keyframes float {
          0%,
          100% {
            transform: translateY(0);
          }
          50% {
            transform: translateY(-20px);
          }
        }
        .animate-float {
          animation: float linear infinite;
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
  );
}


