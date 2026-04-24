declare module "*.vue" {
  import type { DefineComponent } from "vue";

  const component: DefineComponent<object, object, unknown>;
  export default component;
}

declare module "*.css" {
  const content: string;
  export default content;
}

interface Window {
  __SPARK_STATE__?: import("./types").SparkPortalState;
  __SPARK_STATE_PATH__?: string;
  __SPARK_ASSET_BASE_URL__?: string;
  __SPARK_ROUTES__?: Array<{ name: string; method: string; pattern: string }>;
}
