import type { BillingPortalState } from "./types";

const statePath = window.__SPARK_STATE_PATH__ ?? "/billing/state";

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

export function fetchBillingState(): Promise<BillingPortalState> {
  return request<BillingPortalState>(statePath);
}

export function createSubscription(plan: string): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>("/billing/subscription", {
    method: "POST",
    body: JSON.stringify({ plan }),
  });
}

export function updateSubscription(plan: string): Promise<void> {
  return request<void>("/billing/subscription", {
    method: "PUT",
    body: JSON.stringify({ plan }),
  });
}

export function cancelSubscription(): Promise<void> {
  return request<void>("/billing/subscription/cancel", { method: "PUT" });
}

export function resumeSubscription(): Promise<void> {
  return request<void>("/billing/subscription/resume", { method: "PUT" });
}

export function updatePaymentMethod(): Promise<{ transaction_id: string; transaction?: unknown }> {
  return request<{ transaction_id: string; transaction?: unknown }>(
    "/billing/subscription/payment-method",
    {
      method: "PUT",
    },
  );
}

export function markPendingCheckout(checkoutId: string): Promise<void> {
  return request<void>("/billing/pending-checkout", {
    method: "POST",
    body: JSON.stringify({ checkout_id: checkoutId }),
  });
}
