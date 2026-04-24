// Auto-generated from packages/billing route registry. Do not edit by hand.

export type RouteResult = { url: string; method: string };

export function sparkPortal(): RouteResult {
  return { url: "/billing", method: "get" };
}

export function sparkPortalType(params: { type: string | number }): RouteResult {
  return { url: `/billing/${encodeURIComponent(String(params.type))}`, method: "get" };
}

export function sparkPortalBillable(params: {
  type: string | number;
  id: string | number;
}): RouteResult {
  return {
    url: `/billing/${encodeURIComponent(String(params.type))}/${encodeURIComponent(String(params.id))}`,
    method: "get",
  };
}

export function sparkState(): RouteResult {
  return { url: "/billing/state", method: "get" };
}

export function sparkRouteGen(): RouteResult {
  return { url: "/billing/routegen", method: "get" };
}

export function sparkSubscriptionStore(): RouteResult {
  return { url: "/billing/subscription", method: "post" };
}

export function sparkSubscriptionUpdate(): RouteResult {
  return { url: "/billing/subscription", method: "put" };
}

export function sparkSubscriptionCancel(): RouteResult {
  return { url: "/billing/subscription/cancel", method: "put" };
}

export function sparkSubscriptionResume(): RouteResult {
  return { url: "/billing/subscription/resume", method: "put" };
}

export function sparkSubscriptionPaymentMethod(): RouteResult {
  return { url: "/billing/subscription/payment-method", method: "put" };
}

export function sparkPendingCheckout(): RouteResult {
  return { url: "/billing/pending-checkout", method: "post" };
}

export function sparkInvoicesDownload(params: {
  type: string | number;
  id: string | number;
  transaction: string | number;
}): RouteResult {
  return {
    url: `/billing/${encodeURIComponent(String(params.type))}/${encodeURIComponent(String(params.id))}/invoices/${encodeURIComponent(String(params.transaction))}/download`,
    method: "get",
  };
}
