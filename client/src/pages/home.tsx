import { Fragment, useEffect } from "react"
import { StaticMetadata } from "@/contexts/metadata"
import FuckiCloud from "@/components/layout/under-development"
// import AuthFlow from "@/components/layout/authenticate"

export default function Home() {
  // const [signedIn, setSignedIn] = useState(false)
  //
  // useLayoutEffect(() => {
  //   (async () => {
  //     const wasSignedIn = await fetch("/api/auth/status")
  //
  //     if (wasSignedIn.ok) {
  //       setSignedIn(true)
  //     }
  //   })()
  // }, [])

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

      <FuckiCloud />

      {/* {signedIn === false ? ( */}
      {/*   <AuthFlow onSuccess={() => setSignedIn(true)} /> */}
      {/* ) : ( */}
      {/*   <p>Signed In</p> */}
      {/* )} */}
    </Fragment>
  )
}
