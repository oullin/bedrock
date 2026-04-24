import { describe, expect, it } from "vitest";
import { planPrice, statusLabel } from "./format";

describe("Billing portal formatting", () => {
  it("formats minor-unit plan prices", () => {
    expect(
      planPrice({ id: "pri", name: "Pro", interval: "monthly", price: 2900, currency: "USD" }),
    ).toBe("$29.00");
  });

  it("keeps provider-formatted prices unchanged", () => {
    expect(planPrice({ id: "pri", name: "Pro", interval: "monthly", price: "USD 29" })).toBe(
      "USD 29",
    );
  });

  it("maps subscription state labels", () => {
    expect(statusLabel("pending")).toBe("Pending checkout");
    expect(statusLabel("active")).toBe("Active");
    expect(statusLabel("none")).toBe("No subscription");
  });
});
