import { Component, type ReactNode, lazy, Suspense } from "react";

const ClientError = lazy(() => import("@/components/layout/client-error"))

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: any) {
    console.error("Uncaught error:", error, errorInfo);
  }

  render() {
    const { hasError, error } = this.state;
    const { fallback } = this.props;

    if (hasError) {
      return (
        fallback || <Suspense><ClientError resetErrorBoundary={null} error={error || new Error("Unexpected client error, error message was unable to be captured!")} /></Suspense>
      );
    }

    return this.props.children;
  }
}
