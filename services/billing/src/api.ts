import type { BillingPortalState } from "./types";
import {
  billingPendingCheckout,
  billingState,
  billingSubscriptionCancel,
  billingSubscriptionPaymentMethod,
  billingSubscriptionResume,
  billingSubscriptionStore,
  billingSubscriptionUpdate,
  type RouteResult,
} from "./generated/routes";

const statePath = window.__BILLING_STATE_PATH__ ?? billingState().url;

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

export function fetchBillingState(): Promise<BillingPortalState> {
  return request<BillingPortalState>(statePath);
}

export function createSubscription(plan: string): Promise<unknown> {
  return request<unknown>(billingSubscriptionStore(), {
    body: JSON.stringify({ plan }),
  });
}

export function updateSubscription(plan: string): Promise<void> {
  return request<void>(billingSubscriptionUpdate(), {
    body: JSON.stringify({ plan }),
  });
}

export function cancelSubscription(): Promise<void> {
  return request<void>(billingSubscriptionCancel());
}

export function resumeSubscription(): Promise<void> {
  return request<void>(billingSubscriptionResume());
}

export function updatePaymentMethod(): Promise<{ transaction_id: string; transaction?: unknown }> {
  return request<{ transaction_id: string; transaction?: unknown }>(
    billingSubscriptionPaymentMethod(),
  );
}

export function markPendingCheckout(checkoutId: string): Promise<void> {
  return request<void>(billingPendingCheckout(), {
    body: JSON.stringify({ checkout_id: checkoutId }),
  });
}
