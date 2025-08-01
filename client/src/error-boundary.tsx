import { Component, type ReactNode, lazy } from "react";

const ClientError = lazy(() => import("@/pages/client-error"))

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
    const { hasError } = this.state;
    const { fallback } = this.props;

    if (hasError) {
      return (
        fallback || <ClientError />
      );
    }

    return this.props.children;
  }
}
