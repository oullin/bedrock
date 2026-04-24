export interface SparkPlan {
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

export interface SparkInvoice {
  id: string;
  total: string;
  billed_at?: string;
  invoice_url: string;
}

export interface SparkPayment {
  amount: string;
  currency: string;
  date: string;
}

export interface SparkSubscriptionState {
  status: string;
  plan_code: string;
  plan_name: string;
  pending_expires_at: string | null;
  payment_ready_at: string | null;
  portal_url: string;
  pay_now: boolean;
}

export interface SparkCTA {
  visible: boolean;
  label: string;
  remaining_days: number;
}

export interface SparkPortalState {
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
  invoices: SparkInvoice[];
  lastPayment?: SparkPayment | null;
  message?: string;
  monthlyPlans: SparkPlan[];
  nextPayment?: SparkPayment | null;
  paddleSellerId?: number;
  plan?: SparkPlan | null;
  pwAuth?: string;
  pwCustomer?: string | null;
  seatName?: string;
  sparkPath: string;
  state: "none" | "pending" | "active" | "past_due" | "onGracePeriod" | string;
  subscription: SparkSubscriptionState;
  cta: SparkCTA;
  termsUrl?: string;
  yearlyPlans: SparkPlan[];
}
