export type JobStatus = "pending" | "running" | "completed" | "failed" | "silenced" | string;

export interface Job {
  id: string;
  name: string;
  queue: string;
  type: string;
  tags: string[];
  status: JobStatus;
  pushedAt: string | null;
  reservedAt: string | null;
  completedAt: string | null;
  failedAt: string | null;
}

export interface Stats {
  jobsPerMinute: number;
  processes: number;
  queueWithMaxRuntime: string;
  queueWithMaxThroughput: string;
  failedJobs: number;
  pendingJobs: number;
  status: "running" | "paused" | "inactive" | string;
  wait: { longestWait: number };
}

export interface Supervisor {
  name: string;
  pid: number;
  status: string;
  supervisors: string[];
}

export interface MonitoringEntry {
  tag: string;
  count: number;
}

export interface Batch {
  id: string;
  name: string;
  totalJobs: number;
  pendingJobs: number;
  failedJobs: number;
  processedJobs: number;
  createdAt: string;
  finishedAt?: string | null;
  cancelled: boolean;
}

export interface BatchListResponse {
  batches: Batch[];
  nextCursor: string;
}

export interface MetricRow {
  name: string;
  throughput: number;
  averageRuntime: number;
}
