import { describe, expect, it } from "vitest";

import {
  Aggregator,
  Calculator,
  CurrencyManager,
  EUR,
  Exchange,
  GBP,
  JPY,
  MAX_INT64,
  MIN_INT64,
  MoneyConverter,
  MoneyManager,
  Parser,
  SGD,
  USD,
  errors,
  fromSGD,
  fromUSD,
  getDBMoneyValueSeparator,
  moneyFromJSON,
  safeMultiply,
  scanMoney,
  setDBMoneyValueSeparator,
} from "./index";

describe("calculator", () => {
  const calculator = new Calculator();

  it("handles safe arithmetic and overflow-compatible zero results", () => {
    expect(calculator.add(100n, 50n)).toBe(150n);
    expect(calculator.add(MAX_INT64, 1n)).toBe(0n);
    expect(calculator.subtract(100n, 50n)).toBe(50n);
    expect(calculator.subtract(MIN_INT64, 1n)).toBe(0n);
    expect(calculator.multiply(100n, -2n)).toBe(-200n);
    expect(calculator.multiply(MAX_INT64, 2n)).toBe(0n);
  });

  it("divides, allocates, rounds, and rejects overflowing safe multiplication", () => {
    expect(calculator.divide(100n, 3n)).toBe(33n);
    expect(calculator.modulus(100n, 3n)).toBe(1n);
    expect(calculator.allocate(100n, 1n, 3n)).toBe(33n);
    expect(calculator.round(155n, 2)).toBe(200n);
    expect(calculator.round(150n, 2)).toBe(100n);
    expect(() => safeMultiply(MAX_INT64, 2n)).toThrow(errors.overflow);
  });
});

describe("currencies and formatting", () => {
  it("resolves currencies and formats minor units", () => {
    const currencies = new CurrencyManager();
    const sgd = currencies.resolve(SGD);
    const eur = currencies.resolve(EUR);

    expect(sgd.code).toBe(SGD);
    expect(currencies.resolve("missing").code).toBe(SGD);
    expect(currencies.findByNumericCode("978")?.code).toBe(EUR);
    expect(sgd.formatter().format(123456n)).toBe("S$1,234.56");
    expect(eur.formatter().format(-123456n)).toBe("-€1,234.56");
  });
});

describe("parser", () => {
  const parser = new Parser();

  it("detects symbols, ISO codes, defaults, and mixed separators", () => {
    expect(parser.parseAmount("S$1,234.56")).toEqual({ amount: 1234.56, currency: SGD });
    expect(parser.parseAmount("1.234,56 EUR")).toEqual({ amount: 1234.56, currency: EUR });
    expect(parser.parseAmount("10.00", GBP)).toEqual({ amount: 10, currency: GBP });
    expect(parser.parseAmount("C$150.00")).toEqual({ amount: 150, currency: "CAD" });
    expect(parser.parseDecimal("10,50")).toBe(1050);
    expect(parser.parseDecimalWithComma("10,50")).toBe(10.5);
  });

  it("rejects invalid parser inputs", () => {
    expect(() => parser.parseAmount("100.00")).toThrow(errors.currencyNotSpecified);
    expect(() => parser.parseDecimal("1.2,3")).toThrow(errors.invalidMoneyString);
    expect(() => parser.parseAmountString("12.345", 2, false)).toThrow(
      errors.invalidAmountFraction,
    );
  });
});

describe("money manager and money values", () => {
  const manager = new MoneyManager();

  it("creates money from integers, floats, and exact strings", () => {
    expect(manager.create(1500n, SGD).amount()).toBe(1500n);
    expect(manager.createFromFloat(12.349, SGD).amount()).toBe(1235n);
    expect(manager.createFromFloat(98.75, JPY).amount()).toBe(99n);
    expect(manager.createFromString("12.34", SGD).amount()).toBe(1234n);
    expect(manager.createFromString("-.50", SGD).amount()).toBe(-50n);
    expect(manager.createFromString("12.345", "BHD").amount()).toBe(12345n);
    expect(() => manager.createFromString("12.345", SGD)).toThrow(errors.invalidAmountFraction);
  });

  it("compares, displays, and computes money", () => {
    const left = manager.create(200n, SGD);
    const right = manager.create(150n, SGD);
    const euro = manager.create(100n, EUR);

    expect(left.greaterThan(right)).toBe(true);
    expect(right.lessThanOrEqual(left)).toBe(true);
    expect(left.display()).toBe("S$2.00");
    expect(left.asMajorUnits()).toBe(2);
    expect(() => left.assertSameCurrency(euro)).toThrow(errors.currencyMismatch);
    expect(manager.add(left, right).amount()).toBe(350n);
    expect(manager.subtract(left, right).amount()).toBe(50n);
    expect(manager.multiply(right, 2n, 3n).amount()).toBe(900n);
  });

  it("splits and allocates without losing minor units", () => {
    expect(manager.split(manager.create(10n, SGD), 3).map((money) => money.amount())).toEqual([
      4n,
      3n,
      3n,
    ]);
    expect(manager.split(manager.create(-10n, SGD), 3).map((money) => money.amount())).toEqual([
      -4n,
      -3n,
      -3n,
    ]);
    expect(
      manager.allocate(manager.create(101n, SGD), 1, 1, 1).map((money) => money.amount()),
    ).toEqual([34n, 34n, 33n]);
    expect(() => manager.split(manager.create(10n, SGD), 0)).toThrow(errors.invalidSplit);
    expect(() => manager.allocate(manager.create(10n, SGD))).toThrow(errors.noRatiosProvided);
    expect(() => manager.allocate(manager.create(10n, SGD), -1)).toThrow(errors.negativeRatios);
  });

  it("provides factory helpers", () => {
    expect(fromSGD(100n).currency().code).toBe(SGD);
    expect(fromUSD(100n).currency().code).toBe(USD);
  });
});

describe("exchange and aggregation", () => {
  const manager = new MoneyManager();

  it("converts with direct, inverse, and explicit rates", () => {
    const exchange = new Exchange();
    exchange.addRate(SGD, USD, 0.75);
    const converter = new MoneyConverter(manager.getCurrencyManager(), exchange);

    expect(exchange.getRate(USD, SGD)).toBeCloseTo(1 / 0.75);
    expect(converter.convert(manager.create(10000n, SGD), USD).amount()).toBe(7500n);
    expect(converter.convertWithRate(manager.create(10000n, SGD), JPY, 110).amount()).toBe(11_000n);
    expect(() => exchange.getRate(SGD, EUR)).toThrow(errors.currencyConversionNotFound);
  });

  it("aggregates same-currency money", () => {
    const aggregator = new Aggregator(manager);
    const values = [manager.create(100n, SGD), manager.create(50n, SGD), manager.create(25n, SGD)];

    expect(aggregator.sum(...values).amount()).toBe(175n);
    expect(aggregator.min(...values).amount()).toBe(25n);
    expect(aggregator.max(...values).amount()).toBe(100n);
    expect(aggregator.avg(...values).amount()).toBe(58n);
    expect(() => aggregator.sum()).toThrow(errors.noMoneyProvided);
  });
});

describe("serialization", () => {
  const manager = new MoneyManager();

  it("serializes JSON and DB-style values", () => {
    const money = manager.create(1234n, SGD);

    expect(JSON.stringify(money)).toBe('{"amount":1234,"currency":"SGD"}');
    expect(moneyFromJSON('{"amount":12.5,"currency":"SGD"}').amount()).toBe(13n);
    expect(money.value()).toBe("1234|SGD");
    expect(scanMoney("1234|SGD").amount()).toBe(1234n);

    setDBMoneyValueSeparator("~");
    expect(getDBMoneyValueSeparator()).toBe("~");
    expect(manager.create(99n, USD).value()).toBe("99~USD");
    setDBMoneyValueSeparator("|");
  });
});
