export interface BillingPlan {
  id: string;
  name: string;
  interval: string;
  price?: number | string;
  currency?: string;
  features?: string[];
  short_description?: string;
  active?: boolean;
  options?: Record<string, unknown>;
}

export interface BillingInvoice {
  id: string;
  total: string;
  billed_at?: string;
  invoice_url: string;
}

export interface BillingPayment {
  amount: string;
  currency: string;
  date: string;
}

export interface BillingSubscriptionState {
  status: string;
  plan_code: string;
  plan_name: string;
  pending_expires_at: string | null;
  payment_ready_at: string | null;
  portal_url: string;
  pay_now: boolean;
}

export interface BillingCTA {
  visible: boolean;
  label: string;
  remaining_days: number;
}

export interface BillingPortalState {
  appLogo?: string;
  appName: string;
  sandbox: boolean;
  billableId: number | string;
  billableName: string;
  billableType: string;
  brandColor: string;
  clientSideToken?: string;
  dashboardUrl: string;
  defaultInterval: "monthly" | "yearly" | string;
  genericTrialEndsAt?: string | null;
  invoices: BillingInvoice[];
  lastPayment?: BillingPayment | null;
  message?: string;
  monthlyPlans: BillingPlan[];
  nextPayment?: BillingPayment | null;
  paddleSellerId?: number;
  plan?: BillingPlan | null;
  pwAuth?: string;
  pwCustomer?: string | null;
  seatName?: string;
  sparkPath: string;
  state: "none" | "pending" | "active" | "past_due" | "onGracePeriod" | string;
  subscription: BillingSubscriptionState;
  cta: BillingCTA;
  termsUrl?: string;
  yearlyPlans: BillingPlan[];
}
