import React, {
  type ReactElement,
  type ReactNode,
  isValidElement,
  createContext,
  useContext,
  useEffect
} from "react";
import type { MetaTag, LinkTag, ScriptTag, StaticRoute } from "@/types/static.route";
import { useConfig } from "@/contexts/config";
import { useLocation } from "react-router-dom";

const MetadataContext = createContext<boolean | undefined>(undefined);

// Symbol to mark allowed components
const ALLOWED_SYMBOL = Symbol("MetadataAllowed");

type MetadataChild = ReactElement & { type: { [ALLOWED_SYMBOL]?: true } };

// List of internal helper to validate children
function validateChildren(children: ReactNode) {
  React.Children.forEach(children, (child) => {
    if (!isValidElement(child)) return;

    const isAllowed = (child as MetadataChild).type[ALLOWED_SYMBOL];
    if (!isAllowed) {
      const name =
        typeof child.type === "string"
          ? `<${child.type}>`
          : child.type.name
            ? `<${child.type.name}>`
            : "Unknown component";

      throw new Error(
        `${name} is not allowed inside <Metadata>. Use only <Title>, <Meta>, <Link>, or <Script>.`
      );
    }
  });
}

// Main wrapper
export function Metadata({ children }: { children: ReactNode }) {
  validateChildren(children);
  return typeof document !== "undefined" && (
    <MetadataContext.Provider value={true}>
      {children}
    </MetadataContext.Provider>
  );
}

// Hook for inner components
function useMetadataContext(component: string) {
  const context = useContext(MetadataContext);
  if (!context) {
    throw new Error(`<${component}> must be used within <Metadata>`);
  }
}

// Add marker helper
function allowInsideMetadata<T extends React.FC<any>>(component: T): T {
  (component as any)[ALLOWED_SYMBOL] = true;
  return component;
}

export const Title = allowInsideMetadata(({ children }: { children: string }) => {
  useMetadataContext("Title");
  return <title>{children}</title>;
});

export const Meta = allowInsideMetadata((props: MetaTag) => {
  useMetadataContext("Meta");

  if ("name" in props) {
    return <meta name={props.name} content={props.content} />;
  }
  if ("property" in props) {
    return <meta property={props.property} content={props.content} />;
  }
  if ("charset" in props) {
    return <meta charSet={props.charset} />;
  }
  if ("http-equiv" in props) {
    return <meta httpEquiv={props["http-equiv"]} content={props.content} />;
  }
  return null;
});

export const Link = allowInsideMetadata((props: LinkTag) => {
  useMetadataContext("Link");
  return <link {...props} />;
});

type ScriptProps = ScriptTag & React.ScriptHTMLAttributes<HTMLScriptElement>;

export const Script = allowInsideMetadata(({ content, ...rest }: ScriptProps) => {
  useMetadataContext("Script");
  return <script {...rest}>{content}</script>;
});

export const StaticMetadata = () => {
  const { staticRoute } = useConfig();
  const location = useLocation();

  useEffect(() => {
    const serverSentMetadata = document.querySelectorAll('head [id="__SERVER_PROPS__"]');
    serverSentMetadata.forEach((el) => el.remove());
  }, [location]);

  const defaultMeta = staticRoute.find(route => route.path === "*");
  const currentMeta =
    staticRoute.find(route => route.path === location.pathname) ??
    staticRoute.find(route => route.path === "#not_found");

  const mergedMeta = {
    title: currentMeta?.title,
    meta: [...(defaultMeta?.meta ?? []), ...(currentMeta?.meta ?? [])],
    link: [...(defaultMeta?.link ?? []), ...(currentMeta?.link ?? [])],
  };

  return (
    <Metadata>
      {mergedMeta.title && <Title>{mergedMeta.title}</Title>}

      {mergedMeta.meta?.map((meta, i) => {
        if ("charset" in meta) return <Meta key={i} charset={meta.charset} />;
        if ("name" in meta) return <Meta key={i} name={meta.name} content={meta.content} />;
        if ("property" in meta) return <Meta key={i} property={meta.property} content={meta.content} />;
        if ("http-equiv" in meta) return <Meta key={i} http-equiv={meta["http-equiv"]} content={meta.content} />;
        return null;
      })}

      {mergedMeta.link?.map((link, i) => (
        <Link key={i} {...link} />
      ))}
    </Metadata>
  );
};

export const DynamicMetadata = ({ currentMeta }: { currentMeta: StaticRoute }) => {
  const { staticRoute } = useConfig();
  const location = useLocation();

  useEffect(() => {
    const serverSentMetadata = document.querySelectorAll('head [id="__SERVER_PROPS__"]');
    serverSentMetadata.forEach((el) => el.remove());
  }, [location]);

  const defaultMeta = staticRoute.find(route => route.path === "*");

  const mergedMeta = {
    title: currentMeta?.title,
    meta: [...(defaultMeta?.meta ?? []), ...(currentMeta?.meta ?? [])],
    link: [...(defaultMeta?.link ?? []), ...(currentMeta?.link ?? [])],
  };

  return (
    <Metadata>
      {mergedMeta.title && <Title>{mergedMeta.title}</Title>}

      {mergedMeta.meta?.map((meta, i) => {
        if ("charset" in meta) return <Meta key={i} charset={meta.charset} />;
        if ("name" in meta) return <Meta key={i} name={meta.name} content={meta.content} />;
        if ("property" in meta) return <Meta key={i} property={meta.property} content={meta.content} />;
        if ("http-equiv" in meta) return <Meta key={i} http-equiv={meta["http-equiv"]} content={meta.content} />;
        return null;
      })}

      {mergedMeta.link?.map((link, i) => (
        <Link key={i} {...link} />
      ))}
    </Metadata>
  );
};

export const PathBasedMetadata = ({ paths }: { paths: Array<string> }) => {
  const { staticRoute } = useConfig();
  const location = useLocation();

  useEffect(() => {
    const serverSentMetadata = document.querySelectorAll('head [id="__SERVER_PROPS__"]');
    serverSentMetadata.forEach((el) => el.remove());
  }, [location]);

  // Only match explicitly passed paths — no default fallback
  const matchedRoutes = paths
    .map(path => staticRoute.find(route => route.path === path))
    .filter(Boolean);

  const mergedMeta = {
    title: matchedRoutes.find(r => r?.title)?.title,
    meta: matchedRoutes.flatMap(r => r?.meta ?? []),
    link: matchedRoutes.flatMap(r => r?.link ?? []),
  };

  return (
    <Metadata>
      {mergedMeta.title && <Title>{mergedMeta.title}</Title>}

      {mergedMeta.meta.map((meta, i) => {
        if ("charset" in meta) return <Meta key={i} charset={meta.charset} />;
        if ("name" in meta) return <Meta key={i} name={meta.name} content={meta.content} />;
        if ("property" in meta) return <Meta key={i} property={meta.property} content={meta.content} />;
        if ("http-equiv" in meta) return <Meta key={i} http-equiv={meta["http-equiv"]} content={meta.content} />;
        return null;
      })}

      {mergedMeta.link.map((link, i) => (
        <Link key={i} {...link} />
      ))}
    </Metadata>
  );
};
