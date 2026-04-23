export const ROUTES = {
  dashboard: "/",
  monitoring: "/monitoring",
  metrics: "/metrics",
  batches: "/batches",
  jobs: {
    pending: "/jobs/pending",
    completed: "/jobs/completed",
    failed: "/jobs/failed",
    silenced: "/jobs/silenced",
  },
} as const;

type JobRoutes = (typeof ROUTES.jobs)[keyof typeof ROUTES.jobs];
type TopLevelRoutes = Exclude<(typeof ROUTES)[keyof typeof ROUTES], object>;
export type JobQueueRoutePath = TopLevelRoutes | JobRoutes;

export function route(path: JobQueueRoutePath): string {
  return `#${path}`;
}
