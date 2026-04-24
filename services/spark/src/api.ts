import type { SparkPortalState } from "./types";
import {
  sparkPendingCheckout,
  sparkState,
  sparkSubscriptionCancel,
  sparkSubscriptionPaymentMethod,
  sparkSubscriptionResume,
  sparkSubscriptionStore,
  sparkSubscriptionUpdate,
  type RouteResult,
} from "./generated/routes";

const statePath = window.__SPARK_STATE_PATH__ ?? sparkState().url;

async function request<T>(route: RouteResult | string, init: RequestInit = {}): Promise<T> {
  const url = typeof route === "string" ? route : route.url;
  const method = typeof route === "string" ? init.method : route.method.toUpperCase();

  const response = await fetch(url, {
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      ...init.headers,
    },
    method,
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
  return request<Record<string, unknown>>(sparkSubscriptionStore(), {
    body: JSON.stringify({ plan }),
  });
}

export function updateSubscription(plan: string): Promise<void> {
  return request<void>(sparkSubscriptionUpdate(), {
    body: JSON.stringify({ plan }),
  });
}

export function cancelSubscription(): Promise<void> {
  return request<void>(sparkSubscriptionCancel());
}

export function resumeSubscription(): Promise<void> {
  return request<void>(sparkSubscriptionResume());
}

export function updatePaymentMethod(): Promise<{ transaction_id: string; transaction?: unknown }> {
  return request<{ transaction_id: string; transaction?: unknown }>(
    sparkSubscriptionPaymentMethod(),
  );
}

export function markPendingCheckout(checkoutId: string): Promise<void> {
  return request<void>(sparkPendingCheckout(), {
    body: JSON.stringify({ checkout_id: checkoutId }),
  });
}
