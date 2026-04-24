import type { SparkPlan, SparkPortalState } from "./types";

export function statusLabel(state: SparkPortalState["state"]): string {
  switch (state) {
    case "pending":
      return "Pending checkout";
    case "active":
      return "Active";
    case "past_due":
      return "Past due";
    case "onGracePeriod":
      return "Grace period";
    default:
      return "No subscription";
  }
}

export function planPrice(plan: SparkPlan): string {
  if (plan.price == null || plan.price === "") {
    return "Price pending";
  }

  if (typeof plan.price === "number") {
    const amount = plan.price / 100;

    return new Intl.NumberFormat(undefined, {
      style: "currency",
      currency: plan.currency || "USD",
    }).format(amount);
  }

  return String(plan.price);
}
