// Auto-generated from packages/billing route registry. Do not edit by hand.

export type RouteResult = { url: string; method: string };

export function billingPortal(): RouteResult {
  return { url: "/billing", method: "get" };
}

export function billingPortalType(params: { type: string | number }): RouteResult {
  return { url: `/billing/${encodeURIComponent(String(params.type))}`, method: "get" };
}

export function billingPortalBillable(params: {
  type: string | number;
  id: string | number;
}): RouteResult {
  return {
    url: `/billing/${encodeURIComponent(String(params.type))}/${encodeURIComponent(String(params.id))}`,
    method: "get",
  };
}

export function billingState(): RouteResult {
  return { url: "/billing/state", method: "get" };
}

export function billingRouteGen(): RouteResult {
  return { url: "/billing/routegen", method: "get" };
}

export function billingSubscriptionStore(): RouteResult {
  return { url: "/billing/subscription", method: "post" };
}

export function billingSubscriptionUpdate(): RouteResult {
  return { url: "/billing/subscription", method: "put" };
}

export function billingSubscriptionCancel(): RouteResult {
  return { url: "/billing/subscription/cancel", method: "put" };
}

export function billingSubscriptionResume(): RouteResult {
  return { url: "/billing/subscription/resume", method: "put" };
}

export function billingSubscriptionPaymentMethod(): RouteResult {
  return { url: "/billing/subscription/payment-method", method: "put" };
}

export function billingPendingCheckout(): RouteResult {
  return { url: "/billing/pending-checkout", method: "post" };
}

export function billingInvoicesDownload(params: {
  type: string | number;
  id: string | number;
  transaction: string | number;
}): RouteResult {
  return {
    url: `/billing/${encodeURIComponent(String(params.type))}/${encodeURIComponent(String(params.id))}/invoices/${encodeURIComponent(String(params.transaction))}/download`,
    method: "get",
  };
}
