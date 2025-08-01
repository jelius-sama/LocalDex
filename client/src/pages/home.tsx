import { Fragment, useEffect } from "react"
import { StaticMetadata } from "@/contexts/metadata"
import UnderDevelopment from "@/components/layout/under-development"

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

      <UnderDevelopment reminder="Work on Authentication" />
    </Fragment>
  )
}
