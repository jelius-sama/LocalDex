export type ComponentStatus =
    | "ok"
    | "locked"
    | "heavy_load"
    | "degraded";

export interface Healthz {
    status: ComponentStatus;
    timestamp: string; // ISO 8601 UTC string
    components: {
        database: ComponentStatus;
        load: ComponentStatus;
        webserver: ComponentStatus;
        api: ComponentStatus;
    };
}
