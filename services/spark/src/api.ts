import type { SparkPortalState } from "./types";

const statePath = window.__SPARK_STATE_PATH__ ?? "/spark/state";

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(path, {
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      ...init.headers,
    },
    ...init,
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed with ${response.status}`);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

export function fetchSparkState(): Promise<SparkPortalState> {
  return request<SparkPortalState>(statePath);
}

export function createSubscription(plan: string): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>("/spark/subscription", {
    method: "POST",
    body: JSON.stringify({ plan }),
  });
}

export function updateSubscription(plan: string): Promise<void> {
  return request<void>("/spark/subscription", {
    method: "PUT",
    body: JSON.stringify({ plan }),
  });
}

export function cancelSubscription(): Promise<void> {
  return request<void>("/spark/subscription/cancel", { method: "PUT" });
}

export function resumeSubscription(): Promise<void> {
  return request<void>("/spark/subscription/resume", { method: "PUT" });
}

export function updatePaymentMethod(): Promise<{ transaction_id: string; transaction?: unknown }> {
  return request<{ transaction_id: string; transaction?: unknown }>(
    "/spark/subscription/payment-method",
    {
      method: "PUT",
    },
  );
}

export function markPendingCheckout(checkoutId: string): Promise<void> {
  return request<void>("/spark/pending-checkout", {
    method: "POST",
    body: JSON.stringify({ checkout_id: checkoutId }),
  });
}
