import { afterEach, describe, expect, it, vi } from "vitest";
import {
  checkoutOpenOptions,
  handlePaddleCheckoutEvent,
  initializePaddle,
  openPaddleCheckout,
  type CheckoutMode,
  type PaddleGlobal,
} from "./checkout";
import type { SparkPortalState } from "./types";

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("checkoutOpenOptions", () => {
  it("normalizes subscription checkout items for Paddle", () => {
    const options = checkoutOpenOptions({
      customer: "ctm_123",
      custom_data: { source: "spark" },
      items: [{ PriceID: "pri_123", Quantity: 2 }],
    });

    expect(options).toEqual({
      customer: { id: "ctm_123" },
      customData: { source: "spark" },
      items: [{ priceId: "pri_123", quantity: 2 }],
      settings: { allowLogout: false },
    });
  });

  it("opens existing transaction responses by transaction id", () => {
    const options = checkoutOpenOptions({
      transaction: { id: "txn_123" },
    });

    expect(options).toEqual({
      settings: { allowLogout: false },
      transactionId: "txn_123",
    });
  });

  it("opens payment method responses by transaction id", () => {
    const options = checkoutOpenOptions({
      transaction_id: "txn_payment_method",
      transaction: { id: "ignored_nested" },
    });

    expect(options).toEqual({
      settings: { allowLogout: false },
      transactionId: "txn_payment_method",
    });
  });
});

describe("openPaddleCheckout", () => {
  it("launches Paddle checkout", () => {
    const paddle = fakePaddle();
    vi.stubGlobal("window", { Paddle: paddle });

    const options = {
      items: [{ priceId: "pri_123", quantity: 1 }],
      settings: { allowLogout: false },
    };

    openPaddleCheckout(options);

    expect(paddle.Checkout.open).toHaveBeenCalledWith(options);
  });
});

describe("initializePaddle", () => {
  it("configures sandbox and initializes Paddle with checkout callbacks", () => {
    const paddle = fakePaddle();
    vi.stubGlobal("window", { Paddle: paddle });

    const callback = vi.fn();
    const initialized = initializePaddle(
      {
        ...minimalState(),
        clientSideToken: "test_token",
        pwCustomer: "ctm_123",
        sandbox: true,
      },
      callback,
    );

    expect(initialized).toBe(true);
    expect(paddle.Environment?.set).toHaveBeenCalledWith("sandbox");
    expect(paddle.Initialize).toHaveBeenCalledWith({
      eventCallback: callback,
      pwCustomer: "ctm_123",
      token: "test_token",
    });
  });
});

describe("handlePaddleCheckoutEvent", () => {
  it("marks subscription checkouts pending when checkout completes", async () => {
    let mode: CheckoutMode = "subscription";
    const markPendingCheckout = vi.fn(async (_checkoutId: string) => {});
    const refresh = vi.fn(async (_message?: string) => {});
    const setPendingState = vi.fn();

    await handlePaddleCheckoutEvent(
      { name: "checkout.completed", data: { id: "chk_123" } },
      {
        getMode: () => mode,
        markPendingCheckout,
        refresh,
        setMode: (next) => {
          mode = next;
        },
        setPendingState,
      },
    );

    expect(markPendingCheckout).toHaveBeenCalledWith("chk_123");
    expect(setPendingState).toHaveBeenCalled();
    expect(refresh).toHaveBeenCalledWith("Checkout completed. Your subscription is pending.");
    expect(mode).toBeNull();
  });

  it("does not mark payment method checkouts as pending subscription checkouts", async () => {
    let mode: CheckoutMode = "payment-method";
    const markPendingCheckout = vi.fn(async (_checkoutId: string) => {});

    await handlePaddleCheckoutEvent(
      { name: "checkout.completed", data: { id: "chk_payment" } },
      {
        getMode: () => mode,
        markPendingCheckout,
        refresh: vi.fn(async (_message?: string) => {}),
        setMode: (next) => {
          mode = next;
        },
      },
    );

    expect(markPendingCheckout).not.toHaveBeenCalled();
    expect(mode).toBeNull();
  });
});

function fakePaddle(): PaddleGlobal {
  return {
    Checkout: {
      open: vi.fn(),
    },
    Environment: {
      set: vi.fn(),
    },
    Initialize: vi.fn(),
  };
}

function minimalState(): SparkPortalState {
  return {
    appName: "Spark",
    billableId: 10,
    billableName: "Acme",
    billableType: "team",
    brandColor: "bg-gray-800",
    cta: {
      label: "",
      remaining_days: 0,
      visible: false,
    },
    dashboardUrl: "/",
    defaultInterval: "monthly",
    invoices: [],
    monthlyPlans: [],
    sandbox: false,
    sparkPath: "billing",
    state: "none",
    subscription: {
      pay_now: false,
      payment_ready_at: null,
      pending_expires_at: null,
      plan_code: "",
      plan_name: "",
      portal_url: "",
      status: "",
    },
    yearlyPlans: [],
  };
}
