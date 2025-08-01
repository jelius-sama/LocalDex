import { createContext, useContext, useState, type Dispatch, type SetStateAction } from "react"
import AppConfigJSON from "~/client.config.json";
import StaticRouteJSON from "~/static.route.json";
import { type StaticRoute } from "@/types/static.route";

type SSRData = {
  path: string;
  metadata: StaticRoute;
  api_resp: any;
}

type ConfigState = {
  app: typeof AppConfigJSON;
  staticRoute: StaticRoute[];
  ssrData: SSRData | null;
  setSSRData: Dispatch<SetStateAction<SSRData | null>>;
}

const ConfigProviderContext = createContext<ConfigState | undefined>(undefined)

export function ConfigProvider({ children }: { children: React.ReactNode }) {
  const [ssrData, setSSRData] = useState<SSRData | null>(null)

  return (
    <ConfigProviderContext.Provider value={{ app: AppConfigJSON, staticRoute: StaticRouteJSON as StaticRoute[], ssrData: ssrData, setSSRData: setSSRData }}>
      {children}
    </ConfigProviderContext.Provider>
  )
}

export const useConfig = () => {
  const context = useContext(ConfigProviderContext)

  if (context === undefined)
    throw new Error("useConfig must be used within a ConfigProvider")

  return context
}
