import type { SparkPortalState } from "./types";

export type CheckoutMode = "subscription" | "payment-method" | null;

export interface PaddleCheckoutItem {
  priceId: string;
  quantity: number;
}

export interface PaddleCheckoutOpenOptions {
  customer?: Record<string, unknown>;
  customData?: Record<string, unknown>;
  items?: PaddleCheckoutItem[];
  settings: {
    allowLogout: boolean;
  };
  transactionId?: string;
}

export interface PaddleCheckoutEvent {
  name: string;
  data?: {
    id?: string;
  };
}

export interface PaddleInitializeOptions {
  eventCallback: (event: PaddleCheckoutEvent) => void;
  pwCustomer?: string;
  seller?: number;
  token?: string;
}

export interface PaddleGlobal {
  Checkout: {
    open(options: PaddleCheckoutOpenOptions): void;
  };
  Environment?: {
    set(environment: "sandbox"): void;
  };
  Initialize(options: PaddleInitializeOptions): void;
}

interface CheckoutEventHandlers {
  getMode: () => CheckoutMode;
  markPendingCheckout: (checkoutId: string) => Promise<void>;
  refresh: (message?: string) => Promise<void>;
  setMode: (mode: CheckoutMode) => void;
  setPendingState?: () => void;
}

const checkoutSettings = { allowLogout: false };

export function initializePaddle(
  state: SparkPortalState,
  eventCallback: (event: PaddleCheckoutEvent) => void,
): boolean {
  const paddle = window.Paddle;

  if (!paddle) {
    return false;
  }

  if (state.sandbox) {
    paddle.Environment?.set("sandbox");
  }

  const options: PaddleInitializeOptions = { eventCallback };

  if (state.clientSideToken) {
    options.token = state.clientSideToken;
  } else if (state.paddleSellerId) {
    options.seller = state.paddleSellerId;
  }

  if (state.pwCustomer) {
    options.pwCustomer = state.pwCustomer;
  }

  paddle.Initialize(options);

  return true;
}

export function checkoutOpenOptions(response: unknown): PaddleCheckoutOpenOptions {
  const root = record(response) ?? {};
  const payload = record(root.checkout) ?? root;
  const transactionId = transactionID(payload) ?? transactionID(root);

  if (transactionId) {
    return {
      settings: checkoutSettings,
      transactionId,
    };
  }

  const items = normalizeItems(payload.items);

  if (items.length === 0) {
    throw new Error("Checkout response did not include checkout items.");
  }

  const options: PaddleCheckoutOpenOptions = {
    items,
    settings: checkoutSettings,
  };

  const customer = normalizeCustomer(payload.customer);

  if (customer) {
    options.customer = customer;
  }

  const customData = record(payload.customData) ?? record(payload.custom_data);

  if (customData) {
    options.customData = customData;
  }

  return options;
}

export function openPaddleCheckout(options: PaddleCheckoutOpenOptions): void {
  const checkout = window.Paddle?.Checkout;

  if (!checkout) {
    throw new Error("Paddle checkout is unavailable.");
  }

  checkout.open(options);
}

export async function handlePaddleCheckoutEvent(
  event: PaddleCheckoutEvent,
  handlers: CheckoutEventHandlers,
): Promise<void> {
  if (event.name === "checkout.closed") {
    handlers.setMode(null);
    return;
  }

  if (event.name !== "checkout.completed") {
    return;
  }

  const mode = handlers.getMode();

  try {
    if (mode === "payment-method") {
      return;
    }

    if (mode !== "subscription") {
      return;
    }

    if (event.data?.id) {
      await handlers.markPendingCheckout(event.data.id);
      handlers.setPendingState?.();
      await handlers.refresh("Checkout completed. Your subscription is pending.");
      return;
    }

    await handlers.refresh("Checkout completed.");
  } finally {
    handlers.setMode(null);
  }
}

function normalizeItems(items: unknown): PaddleCheckoutItem[] {
  if (!Array.isArray(items)) {
    return [];
  }

  return items
    .map((item) => {
      const row = record(item);
      const priceId = stringValue(row?.priceId) ?? stringValue(row?.PriceID);

      if (!priceId) {
        return null;
      }

      return {
        priceId,
        quantity: numberValue(row?.quantity) ?? numberValue(row?.Quantity) ?? 1,
      };
    })
    .filter((item): item is PaddleCheckoutItem => item !== null);
}

function normalizeCustomer(customer: unknown): Record<string, unknown> | undefined {
  if (typeof customer === "string" && customer !== "") {
    return { id: customer };
  }

  return record(customer);
}

function transactionID(payload: Record<string, unknown> | undefined): string | undefined {
  if (!payload) {
    return undefined;
  }

  const transaction = record(payload.transaction);

  return (
    stringValue(payload.transactionId) ??
    stringValue(payload.transaction_id) ??
    stringValue(transaction?.id) ??
    stringValue(transaction?.transaction_id)
  );
}

function record(value: unknown): Record<string, unknown> | undefined {
  if (value === null || typeof value !== "object" || Array.isArray(value)) {
    return undefined;
  }

  return value as Record<string, unknown>;
}

function stringValue(value: unknown): string | undefined {
  return typeof value === "string" && value !== "" ? value : undefined;
}

function numberValue(value: unknown): number | undefined {
  return typeof value === "number" && Number.isFinite(value) ? value : undefined;
}
