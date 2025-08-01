import { Fragment, useEffect } from "react"
import { StaticMetadata } from "@/contexts/metadata"

export default function Home() {
  useEffect(() => {
    const handleImageLoad = () => {
      const event = new CustomEvent("PageLoaded", {
        detail: { pathname: window.location.pathname },
      });
      window.dispatchEvent(event);
    };

    handleImageLoad();
  }, []);

  return (
    <Fragment>
      <StaticMetadata />

      <h1>Hello</h1>
    </Fragment>
  )
}
