export * from "./currency-data";

import {
  SGD,
  currencyDefinitions,
  currencySymbols,
  type CurrencyCode,
  type CurrencyDefinition,
  type CurrencySymbol,
} from "./currency-data";

export type Amount = bigint;
export type AmountInput = bigint | number | string;

export const MAX_INT64 = 9_223_372_036_854_775_807n;
export const MIN_INT64 = -9_223_372_036_854_775_808n;

export class MoneyError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "MoneyError";
  }
}

export const errors = {
  currencyMismatch: new MoneyError("currencies don't match"),
  currencyNotFound: new MoneyError("currency not found"),
  currencyConversionNotFound: new MoneyError("currency conversion rate not found"),
  noCurrencyInstance: new MoneyError("money instance has no currency"),
  noCurrencyManager: new MoneyError("currency manager cannot be nil"),
  noCurrencyMapDataset: new MoneyError("currency map dataset cannot be nil or empty"),
  invalidJsonUnmarshal: new MoneyError("invalid json unmarshal"),
  invalidExchangeRate: new MoneyError("invalid exchange rate"),
  noJsonParserProvided: new MoneyError("no json parser provided"),
  jsonUnmarshalFuncNil: new MoneyError("money.JSON: unmarshal function cannot be nil"),
  jsonMarshalFuncNil: new MoneyError("money.JSON: marshal function cannot be nil"),
  noMoneyProvided: new MoneyError("no money objects provided"),
  invalidMoneyString: new MoneyError("invalid money string format"),
  currencyNotSpecified: new MoneyError("currency not specified or detected"),
  noMultipliersProvided: new MoneyError("no multipliers provided"),
  noConverterProvided: new MoneyError("no converter provided"),
  emptyAmountString: new MoneyError("amount string cannot be empty"),
  invalidAmountMultiple: new MoneyError("invalid amount: multiple decimal points"),
  invalidAmountFraction: new MoneyError("too many decimal places for curr"),
  invalidAmount: new MoneyError("invalid amount"),
  invalidSplit: new MoneyError("split must be higher than zero"),
  noRatiosProvided: new MoneyError("no ratios specified"),
  negativeRatios: new MoneyError("negative ratios not allowed"),
  ratiosExceedMaxInt: new MoneyError("sum of given ratios exceeds max int"),
  invalidAggregatorProvider: new MoneyError("invalid aggregator: nil manager"),
  overflow: new MoneyError("arithmetic operation resulted in overflow"),
  parserNotProvided: new MoneyError("parser was not provided"),
  parserInvalidState: new MoneyError("parser is nil or iso is nil"),
} as const;

function amountInputToBigInt(value: AmountInput): bigint {
  if (typeof value === "bigint") {
    return value;
  }

  if (typeof value === "number") {
    if (!Number.isSafeInteger(value)) {
      throw errors.invalidAmount;
    }

    return BigInt(value);
  }

  if (!/^[+-]?\d+$/.test(value.trim())) {
    throw errors.invalidAmount;
  }

  return BigInt(value.trim());
}

function ensureInt64(value: bigint): bigint {
  if (value > MAX_INT64 || value < MIN_INT64) {
    throw errors.overflow;
  }

  return value;
}

function safeNumberToBigInt(value: number): bigint {
  if (!Number.isFinite(value) || !Number.isSafeInteger(value)) {
    throw errors.invalidAmount;
  }

  return BigInt(value);
}

function roundHalfAwayFromZero(value: number): bigint {
  return safeNumberToBigInt(value < 0 ? Math.ceil(value - 0.5) : Math.floor(value + 0.5));
}

export class Formatter {
  constructor(
    public readonly fraction: number,
    public readonly decimal: string,
    public readonly thousand: string,
    public readonly grapheme: string,
    public readonly template: string,
  ) {}

  format(amountInput: AmountInput): string {
    const amount = amountInputToBigInt(amountInput);
    let formatted = (amount < 0n ? -amount : amount).toString();

    if (formatted.length <= this.fraction) {
      formatted = "0".repeat(this.fraction - formatted.length + 1) + formatted;
    }

    if (this.thousand !== "") {
      for (let index = formatted.length - this.fraction - 3; index > 0; index -= 3) {
        formatted = `${formatted.slice(0, index)}${this.thousand}${formatted.slice(index)}`;
      }
    }

    if (this.fraction > 0) {
      formatted = `${formatted.slice(0, formatted.length - this.fraction)}${this.decimal}${formatted.slice(formatted.length - this.fraction)}`;
    }

    formatted = this.template.replace("1", formatted).replace("$", this.grapheme);
    return amount < 0n ? `-${formatted}` : formatted;
  }

  toMajorUnits(amountInput: AmountInput): number {
    const amount = amountInputToBigInt(amountInput);
    return this.fraction === 0 ? Number(amount) : Number(amount) / 10 ** this.fraction;
  }
}

export class Currency {
  readonly decimal: string;
  readonly thousand: string;
  readonly code: CurrencyCode;
  readonly fraction: number;
  readonly numericCode: string;
  readonly grapheme: string;
  readonly template: string;

  constructor(definition: CurrencyDefinition) {
    this.decimal = definition.decimal;
    this.thousand = definition.thousand;
    this.code = definition.code;
    this.fraction = definition.fraction;
    this.numericCode = definition.numericCode;
    this.grapheme = definition.grapheme;
    this.template = definition.template;
  }

  formatter(): Formatter {
    return new Formatter(this.fraction, this.decimal, this.thousand, this.grapheme, this.template);
  }

  equals(other: Currency | null | undefined): boolean {
    return other instanceof Currency && this.code === other.code;
  }

  valueOf(): CurrencyCode {
    return this.code;
  }
}

function normalizeCurrencyCode(code: string): string {
  return code.trim().toUpperCase();
}

export class CurrencyMap {
  private readonly dataset: Map<string, Currency>;

  constructor(dataset: Record<string, CurrencyDefinition | Currency> = currencyDefinitions) {
    this.dataset = new Map(
      Object.entries(dataset).map(([code, value]) => [
        normalizeCurrencyCode(code),
        value instanceof Currency ? value : new Currency(value),
      ]),
    );
  }

  get(code: string): Currency | undefined {
    return this.dataset.get(normalizeCurrencyCode(code));
  }

  isEmpty(): boolean {
    return this.dataset.size === 0;
  }

  isNotEmpty(): boolean {
    return !this.isEmpty();
  }

  findByCode(code: string): Currency | undefined {
    return this.get(code);
  }

  getCodes(): string[] {
    return Array.from(this.dataset.keys());
  }

  set(code: string, currency: Currency): void {
    this.dataset.set(normalizeCurrencyCode(code), currency);
  }
}

export function newCurrenciesMap(): CurrencyMap {
  return new CurrencyMap();
}

export interface CurrencyProvider {
  get(): CurrencyDefinition;
  getCode(): string;
  getSymbols(): readonly CurrencySymbol[];
}

export class DefaultProvider implements CurrencyProvider {
  getCode(): CurrencyCode {
    return SGD;
  }

  get(): CurrencyDefinition {
    return currencyDefinitions.SGD;
  }

  getSymbols(): readonly CurrencySymbol[] {
    return currencySymbols;
  }
}

export class CurrencyManager {
  private readonly currencies: CurrencyMap;
  private readonly symbols: readonly CurrencySymbol[];
  private readonly defaultCurrency: Currency;

  constructor(
    provider: CurrencyProvider = new DefaultProvider(),
    dataset: Record<string, CurrencyDefinition | Currency> = currencyDefinitions,
  ) {
    this.symbols = provider.getSymbols();
    this.currencies = new CurrencyMap(dataset);
    this.defaultCurrency = new Currency(provider.get());
  }

  findByCode(code: string): Currency | undefined {
    return this.currencies.findByCode(code);
  }

  findByNumericCode(code: string): Currency | undefined {
    const lookup = normalizeCurrencyCode(code);
    return this.currencies
      .getCodes()
      .map((currencyCode) => this.currencies.get(currencyCode))
      .find((currency) => currency?.numericCode.toUpperCase() === lookup);
  }

  add(currency: Currency | CurrencyDefinition): Currency | undefined {
    const next = currency instanceof Currency ? currency : new Currency(currency);
    const code = normalizeCurrencyCode(next.code);

    if (code === "") {
      return undefined;
    }

    this.currencies.set(code, next);
    return next;
  }

  addFrom(
    code: string,
    grapheme: string,
    template: string,
    decimal: string,
    thousand: string,
    numericCode: string,
    fraction: number,
  ): Currency | undefined {
    return this.add({
      code: normalizeCurrencyCode(code) as CurrencyCode,
      grapheme,
      template,
      decimal,
      thousand,
      numericCode: normalizeCurrencyCode(numericCode),
      fraction,
    });
  }

  getDefault(): Currency {
    return this.defaultCurrency;
  }

  getSymbols(): readonly CurrencySymbol[] {
    return this.symbols;
  }

  resolve(code: string): Currency {
    return this.findByCode(code) ?? this.defaultCurrency;
  }
}

export class ISOCodePattern {
  private readonly currencies: CurrencyMap;
  private readonly symbols: readonly CurrencySymbol[];
  private symbolsLongestFirst?: readonly CurrencySymbol[];
  private pattern?: RegExp;

  constructor(
    provider: CurrencyProvider = new DefaultProvider(),
    dataset: Record<string, CurrencyDefinition | Currency> = currencyDefinitions,
  ) {
    this.symbols = provider.getSymbols();
    this.currencies = new CurrencyMap(dataset);
  }

  getSymbolsLongestFirst(): readonly CurrencySymbol[] {
    this.symbolsLongestFirst ??= [...this.symbols].sort(
      (left, right) => right.id.length - left.id.length || left.id.localeCompare(right.id),
    );
    return this.symbolsLongestFirst;
  }

  getPattern(): RegExp {
    this.pattern ??= new RegExp(`\\b(${this.currencies.getCodes().join("|")})\\b`);
    return this.pattern;
  }
}

export class Calculator {
  add(left: AmountInput, right: AmountInput): bigint {
    return safeAdd(left, right);
  }

  subtract(left: AmountInput, right: AmountInput): bigint {
    return safeSubtract(left, right);
  }

  multiply(amount: AmountInput, seed: AmountInput): bigint {
    return ration(amount, seed);
  }

  safeMultiply(initial: AmountInput, ...multipliers: AmountInput[]): bigint {
    return safeMultiply(initial, ...multipliers);
  }

  divide(amount: AmountInput, seed: AmountInput): bigint {
    const divisor = amountInputToBigInt(seed);
    return divisor === 0n ? 0n : amountInputToBigInt(amount) / divisor;
  }

  modulus(amount: AmountInput, seed: AmountInput): bigint {
    const divisor = amountInputToBigInt(seed);
    return divisor === 0n ? 0n : amountInputToBigInt(amount) % divisor;
  }

  allocate(amount: AmountInput, ratio: AmountInput, scale: AmountInput): bigint {
    const base = amountInputToBigInt(amount);
    const divisor = amountInputToBigInt(scale);

    if (base === 0n || divisor === 0n) {
      return 0n;
    }

    return ration(base, ratio) / divisor;
  }

  absolute(amount: AmountInput): bigint {
    const value = amountInputToBigInt(amount);
    return value < 0n ? -value : value;
  }

  negative(amount: AmountInput): bigint {
    const value = amountInputToBigInt(amount);
    return value > 0n ? -value : value;
  }

  round(amount: AmountInput, exponent: number): bigint {
    const value = amountInputToBigInt(amount);

    if (value === 0n || exponent <= 0 || exponent > 18) {
      return value;
    }

    const unit = 10n ** BigInt(exponent);
    let absolute = value < 0n ? -value : value;
    const remainder = absolute % unit;

    if (remainder > unit / 2n) {
      absolute += unit;
    }

    absolute = (absolute / unit) * unit;
    return value < 0n ? -absolute : absolute;
  }
}

export function safeAdd(left: AmountInput, right: AmountInput): bigint {
  const result = amountInputToBigInt(left) + amountInputToBigInt(right);
  return result > MAX_INT64 || result < MIN_INT64 ? 0n : result;
}

export function safeSubtract(left: AmountInput, right: AmountInput): bigint {
  const result = amountInputToBigInt(left) - amountInputToBigInt(right);
  return result > MAX_INT64 || result < MIN_INT64 ? 0n : result;
}

export function ration(amountInput: AmountInput, ratioInput: AmountInput): bigint {
  const amount = amountInputToBigInt(amountInput);
  const ratio = amountInputToBigInt(ratioInput);

  if (amount === 0n || ratio === 0n) {
    return 0n;
  }

  const result = amount * ratio;
  return result > MAX_INT64 || result < MIN_INT64 ? 0n : result;
}

export function safeMultiply(initial: AmountInput, ...multipliers: AmountInput[]): bigint {
  let result = amountInputToBigInt(initial);

  for (const multiplierInput of multipliers) {
    const multiplier = amountInputToBigInt(multiplierInput);
    const next = result * multiplier;

    if (next > MAX_INT64 || next < MIN_INT64) {
      throw errors.overflow;
    }

    result = next;
  }

  return result;
}

export class Parser {
  constructor(private readonly iso: ISOCodePattern = new ISOCodePattern()) {}

  parseAmount(input: string, defaultCurrency = ""): { amount: number; currency: string } {
    const trimmed = input.trim();

    if (trimmed === "") {
      throw errors.invalidMoneyString;
    }

    const extracted = this.extractCurrency(trimmed, defaultCurrency);
    return {
      amount: this.parseNumericString(extracted.input, false),
      currency: extracted.currency,
    };
  }

  parseAmountWithDecimalComma(
    input: string,
    defaultCurrency = "",
  ): { amount: number; currency: string } {
    const trimmed = input.trim();

    if (trimmed === "") {
      throw errors.invalidMoneyString;
    }

    const extracted = this.extractCurrency(trimmed, defaultCurrency);
    return { amount: this.parseNumericString(extracted.input, true), currency: extracted.currency };
  }

  parseDecimal(amount: string): number {
    return this.parseNumericString(amount, false);
  }

  parseDecimalWithComma(amount: string): number {
    return this.parseNumericString(amount, true);
  }

  parseStringSign(amount: string): { amount: string; negative: boolean } {
    if (amount.startsWith("-")) {
      return { amount: amount.slice(1), negative: true };
    }

    if (amount.startsWith("+")) {
      return { amount: amount.slice(1), negative: false };
    }

    return { amount, negative: false };
  }

  parseDecimalParts(amount: string): { integerPart: string; decimalPart: string } {
    const parts = amount.split(".");

    if (parts.length > 2) {
      throw errors.invalidAmountMultiple;
    }

    return { integerPart: parts[0] === "" ? "0" : parts[0], decimalPart: parts[1] ?? "" };
  }

  validateAndPadDecimal(decimalPart: string, fraction: number): string {
    if (decimalPart.length > fraction) {
      throw errors.invalidAmountFraction;
    }

    return decimalPart.padEnd(fraction, "0");
  }

  parseAmountString(amount: string, fraction: number, negative: boolean): bigint {
    const { integerPart, decimalPart } = this.parseDecimalParts(amount);
    const padded = this.validateAndPadDecimal(decimalPart, fraction);
    const combined = `${integerPart}${padded}`;

    if (!/^\d+$/.test(combined)) {
      throw errors.invalidAmount;
    }

    const value = ensureInt64(BigInt(combined));
    return negative ? -value : value;
  }

  private parseNumericString(input: string, useDecimalComma: boolean): number {
    let normalized = input.trim().replaceAll(" ", "");
    const hasDot = normalized.includes(".");
    const hasComma = normalized.includes(",");

    if (hasDot && hasComma) {
      const lastDot = normalized.lastIndexOf(".");
      const lastComma = normalized.lastIndexOf(",");

      if (lastDot < lastComma) {
        if (!this.validThousandsGrouping(normalized, ".", ",")) {
          throw errors.invalidMoneyString;
        }

        normalized = normalized.replaceAll(".", "").replaceAll(",", ".");
      } else {
        if (!this.validThousandsGrouping(normalized, ",", ".")) {
          throw errors.invalidMoneyString;
        }

        normalized = normalized.replaceAll(",", "");
      }
    } else if (hasComma) {
      normalized = useDecimalComma
        ? normalized.replaceAll(",", ".")
        : normalized.replaceAll(",", "");
    }

    if (!/^[+-]?(?:\d+(?:\.\d+)?|\.\d+)$/.test(normalized)) {
      throw errors.invalidMoneyString;
    }

    return Number.parseFloat(normalized);
  }

  private validThousandsGrouping(input: string, thousandsSep: string, decimalSep: string): boolean {
    const decimalIndex = input.lastIndexOf(decimalSep);
    let integerPart = decimalIndex === -1 ? input : input.slice(0, decimalIndex);

    if (integerPart.startsWith("-") || integerPart.startsWith("+")) {
      integerPart = integerPart.slice(1);
    }

    const groups = integerPart.split(thousandsSep);

    if (groups.length === 1) {
      return true;
    }

    if (groups[0] === "" || groups[0].length > 3) {
      return false;
    }

    return groups.slice(1).every((group) => group.length === 3);
  }

  private extractCurrency(
    input: string,
    defaultCurrency: string,
  ): { input: string; currency: string } {
    let currency = defaultCurrency;
    let nextInput = input;

    for (const symbol of this.iso.getSymbolsLongestFirst()) {
      if (nextInput.includes(symbol.id)) {
        currency = symbol.currency;
        nextInput = nextInput.replaceAll(symbol.id, "");
        break;
      }
    }

    nextInput = nextInput.trim();
    const match = this.iso.getPattern().exec(nextInput);

    if (match?.[1]) {
      currency = match[1];
      nextInput = `${nextInput.slice(0, match.index)}${nextInput.slice(match.index + match[1].length)}`;
    }

    if (currency === "") {
      throw errors.currencyNotSpecified;
    }

    return { input: nextInput, currency };
  }
}

export class Money {
  constructor(
    private readonly amountValue: bigint,
    private readonly currencyValue: Currency,
  ) {}

  amount(): bigint {
    return this.amountValue;
  }
  currency(): Currency {
    return this.currencyValue;
  }

  assertSameCurrency(other: Money | null | undefined): void {
    if (!other) throw errors.noMoneyProvided;
    if (!this.sameCurrency(other)) throw errors.currencyMismatch;
  }

  sameCurrency(other: Money | null | undefined): boolean {
    if (!other) throw errors.noMoneyProvided;
    return this.currencyValue.equals(other.currencyValue);
  }

  compareAmount(other: Money): -1 | 0 | 1 {
    this.assertSameCurrency(other);
    if (this.amountValue > other.amountValue) return 1;
    if (this.amountValue < other.amountValue) return -1;
    return 0;
  }

  equals(other: Money): boolean {
    return this.compareAmount(other) === 0;
  }
  greaterThan(other: Money): boolean {
    return this.compareAmount(other) === 1;
  }
  greaterThanOrEqual(other: Money): boolean {
    return this.compareAmount(other) >= 0;
  }
  lessThan(other: Money): boolean {
    return this.compareAmount(other) === -1;
  }
  lessThanOrEqual(other: Money): boolean {
    return this.compareAmount(other) <= 0;
  }
  isZero(): boolean {
    return this.amountValue === 0n;
  }
  isPositive(): boolean {
    return this.amountValue > 0n;
  }
  isNegative(): boolean {
    return this.amountValue < 0n;
  }
  display(): string {
    return this.currencyValue.formatter().format(this.amountValue);
  }
  asMajorUnits(): number {
    return this.currencyValue.formatter().toMajorUnits(this.amountValue);
  }
  compare(other: Money): -1 | 0 | 1 {
    return this.compareAmount(other);
  }

  toJSON(): { amount: number | string; currency: CurrencyCode } {
    const amount =
      this.amountValue <= BigInt(Number.MAX_SAFE_INTEGER) &&
      this.amountValue >= BigInt(Number.MIN_SAFE_INTEGER)
        ? Number(this.amountValue)
        : this.amountValue.toString();
    return { amount, currency: this.currencyValue.code };
  }

  value(): string {
    return `${this.amountValue.toString()}${getDBMoneyValueSeparator()}${this.currencyValue.code}`;
  }
}

export class MoneyManager {
  private readonly parser = new Parser();
  private readonly calculator = new Calculator();

  constructor(private readonly currencyManager: CurrencyManager = new CurrencyManager()) {}

  static withCurrencyManager(currencyManager: CurrencyManager | null | undefined): MoneyManager {
    if (!currencyManager) throw errors.noCurrencyManager;
    return new MoneyManager(currencyManager);
  }

  create(amount: AmountInput, code: string): Money {
    return new Money(ensureInt64(amountInputToBigInt(amount)), this.currencyManager.resolve(code));
  }

  createFromFloat(amount: number, code: string): Money {
    const currency = this.currencyManager.resolve(code);
    return this.create(roundHalfAwayFromZero(amount * 10 ** currency.fraction), code);
  }

  createFromString(amount: string, code: string): Money {
    const trimmed = amount.trim();
    if (trimmed === "") throw errors.emptyAmountString;
    const currency = this.currencyManager.resolve(code);
    const sign = this.parser.parseStringSign(trimmed);
    const value = this.parser.parseAmountString(sign.amount, currency.fraction, sign.negative);
    return this.create(value, code);
  }

  getCurrencyManager(): CurrencyManager {
    return this.currencyManager;
  }

  add(first: Money, ...rest: Money[]): Money {
    let amount = first.amount();
    for (const item of rest) {
      first.assertSameCurrency(item);
      amount = safeAdd(amount, item.amount());
    }
    return this.create(amount, first.currency().code);
  }

  subtract(first: Money, ...rest: Money[]): Money {
    let amount = first.amount();
    for (const item of rest) {
      first.assertSameCurrency(item);
      amount = safeSubtract(amount, item.amount());
    }
    return this.create(amount, first.currency().code);
  }

  multiply(money: Money, ...values: AmountInput[]): Money {
    if (values.length === 0) throw errors.noMultipliersProvided;
    return this.create(
      this.calculator.safeMultiply(money.amount(), ...values),
      money.currency().code,
    );
  }

  absolute(money: Money): Money {
    return this.create(this.calculator.absolute(money.amount()), money.currency().code);
  }
  negative(money: Money): Money {
    return this.create(-money.amount(), money.currency().code);
  }
  round(money: Money): Money {
    return this.create(
      this.calculator.round(money.amount(), money.currency().fraction),
      money.currency().code,
    );
  }

  split(money: Money, count: number): Money[] {
    if (count <= 0) throw errors.invalidSplit;
    const divisor = BigInt(count);
    const quotient = money.amount() / divisor;
    const remainder = money.amount() % divisor;
    const result = Array.from({ length: count }, () =>
      this.create(quotient, money.currency().code),
    );
    const increment = money.amount() < 0n ? -1n : 1n;
    for (let index = 0; index < Number(remainder < 0n ? -remainder : remainder); index += 1) {
      result[index] = this.create(result[index].amount() + increment, money.currency().code);
    }
    return result;
  }

  allocate(money: Money, ...ratios: number[]): Money[] {
    if (ratios.length === 0) throw errors.noRatiosProvided;
    let ratioTotal = 0n;
    for (const ratio of ratios) {
      if (ratio < 0) throw errors.negativeRatios;
      ratioTotal += BigInt(ratio);
      if (ratioTotal > MAX_INT64) throw errors.ratiosExceedMaxInt;
    }
    const result = ratios.map((ratio) =>
      this.create(
        this.calculator.allocate(money.amount(), BigInt(ratio), ratioTotal),
        money.currency().code,
      ),
    );
    if (ratioTotal === 0n) return result;
    let allocatedTotal = result.reduce((total, item) => total + item.amount(), 0n);
    let leftover = money.amount() - allocatedTotal;
    const increment = leftover < 0n ? -1n : 1n;
    for (let index = 0; leftover !== 0n && index < result.length; index += 1) {
      result[index] = this.create(result[index].amount() + increment, money.currency().code);
      allocatedTotal += increment;
      leftover = money.amount() - allocatedTotal;
    }
    return result;
  }
}

export class Exchange {
  private readonly rates = new Map<string, Map<string, number>>();
  isValid(): boolean {
    return this.rates.size > 0;
  }
  isInvalid(): boolean {
    return !this.isValid();
  }

  addRate(baseCurrency: string, counterCurrency: string, rate: number): void {
    if (rate <= 0) throw errors.invalidExchangeRate;
    const base = normalizeCurrencyCode(baseCurrency);
    const counter = normalizeCurrencyCode(counterCurrency);
    const rates = this.rates.get(base) ?? new Map<string, number>();
    rates.set(counter, rate);
    this.rates.set(base, rates);
  }

  getRate(baseCurrency: string, counterCurrency: string): number {
    const base = normalizeCurrencyCode(baseCurrency);
    const counter = normalizeCurrencyCode(counterCurrency);
    if (base === counter) return 1;
    const direct = this.rates.get(base)?.get(counter);
    if (direct !== undefined) return direct;
    const inverse = this.rates.get(counter)?.get(base);
    if (inverse !== undefined) return 1 / inverse;
    throw errors.currencyConversionNotFound;
  }

  convertAmount(
    amountInput: AmountInput,
    fromCurrencyCode: string,
    fromFraction: number,
    toCurrencyCode: string,
    toFraction: number,
  ): bigint {
    const amount = amountInputToBigInt(amountInput);
    if (normalizeCurrencyCode(fromCurrencyCode) === normalizeCurrencyCode(toCurrencyCode))
      return amount;
    return this.convertAmountWithRate(
      amount,
      fromFraction,
      toFraction,
      this.getRate(fromCurrencyCode, toCurrencyCode),
    );
  }

  convertAmountWithRate(
    amountInput: AmountInput,
    fromFraction: number,
    toFraction: number,
    rate: number,
  ): bigint {
    if (rate <= 0) throw errors.invalidExchangeRate;
    const amount = Number(amountInputToBigInt(amountInput));
    return roundHalfAwayFromZero((amount / 10 ** fromFraction) * rate * 10 ** toFraction);
  }
}

export class MoneyConverter {
  constructor(
    private readonly currencies: CurrencyManager,
    private readonly exchange: Exchange,
  ) {
    if (exchange.isInvalid()) throw errors.invalidExchangeRate;
  }

  convert(money: Money, toCurrency: string): Money {
    const target = this.currencies.findByCode(toCurrency);
    if (!target) throw errors.currencyNotFound;
    return new Money(
      this.exchange.convertAmount(
        money.amount(),
        money.currency().code,
        money.currency().fraction,
        target.code,
        target.fraction,
      ),
      target,
    );
  }

  convertWithRate(money: Money, toCurrency: string, rate: number): Money {
    const target = this.currencies.findByCode(toCurrency);
    if (!target) throw errors.currencyNotFound;
    return new Money(
      this.exchange.convertAmountWithRate(
        money.amount(),
        money.currency().fraction,
        target.fraction,
        rate,
      ),
      target,
    );
  }
}

export class Aggregator {
  constructor(private readonly manager: MoneyManager) {}
  sum(...moneys: Money[]): Money {
    if (moneys.length === 0) throw errors.noMoneyProvided;
    return this.manager.add(moneys[0], ...moneys.slice(1));
  }
  min(...moneys: Money[]): Money {
    if (moneys.length === 0) throw errors.noMoneyProvided;
    return moneys.slice(1).reduce((minimum, current) => {
      minimum.assertSameCurrency(current);
      return current.amount() < minimum.amount() ? current : minimum;
    }, moneys[0]);
  }
  max(...moneys: Money[]): Money {
    if (moneys.length === 0) throw errors.noMoneyProvided;
    return moneys.slice(1).reduce((maximum, current) => {
      maximum.assertSameCurrency(current);
      return current.amount() > maximum.amount() ? current : maximum;
    }, moneys[0]);
  }
  avg(...moneys: Money[]): Money {
    if (moneys.length === 0) throw errors.noMoneyProvided;
    const sum = this.sum(...moneys);
    return this.manager.create(sum.amount() / BigInt(moneys.length), sum.currency().code);
  }
}

let dbMoneyValueSeparator = "|";
export function getDBMoneyValueSeparator(): string {
  return dbMoneyValueSeparator;
}
export function setDBMoneyValueSeparator(separator: string): void {
  if (separator.trim() === "") throw new MoneyError(`separator [${separator}] cannot be empty`);
  dbMoneyValueSeparator = separator;
}

export function scanMoney(value: string | Uint8Array, manager = new MoneyManager()): Money {
  const source = typeof value === "string" ? value : new TextDecoder().decode(value);
  const parts = source.split(getDBMoneyValueSeparator());
  if (parts.length !== 2 || parts[0] === "" || parts[1] === "") throw errors.invalidMoneyString;
  return manager.create(parts[0], parts[1]);
}

export function moneyFromJSON(
  value: string | { amount?: number | string | bigint; currency?: string },
  manager = new MoneyManager(),
): Money {
  const raw =
    typeof value === "string"
      ? (JSON.parse(value) as { amount?: number | string | bigint; currency?: string })
      : value;
  const amount = raw.amount ?? 0;
  const currency =
    raw.currency?.trim() === "" || raw.currency === undefined
      ? manager.getCurrencyManager().getDefault().code
      : raw.currency;
  return manager.create(
    typeof amount === "number" ? roundHalfAwayFromZero(amount) : amount,
    currency,
  );
}

const defaultManager = new MoneyManager();
export function create(amount: AmountInput, code: string): Money {
  return defaultManager.create(amount, code);
}

export function fromAED(amount: AmountInput): Money {
  return defaultManager.create(amount, "AED");
}

export function fromAFN(amount: AmountInput): Money {
  return defaultManager.create(amount, "AFN");
}

export function fromALL(amount: AmountInput): Money {
  return defaultManager.create(amount, "ALL");
}

export function fromAMD(amount: AmountInput): Money {
  return defaultManager.create(amount, "AMD");
}

export function fromANG(amount: AmountInput): Money {
  return defaultManager.create(amount, "ANG");
}

export function fromAOA(amount: AmountInput): Money {
  return defaultManager.create(amount, "AOA");
}

export function fromARS(amount: AmountInput): Money {
  return defaultManager.create(amount, "ARS");
}

export function fromAUD(amount: AmountInput): Money {
  return defaultManager.create(amount, "AUD");
}

export function fromAWG(amount: AmountInput): Money {
  return defaultManager.create(amount, "AWG");
}

export function fromAZN(amount: AmountInput): Money {
  return defaultManager.create(amount, "AZN");
}

export function fromBAM(amount: AmountInput): Money {
  return defaultManager.create(amount, "BAM");
}

export function fromBBD(amount: AmountInput): Money {
  return defaultManager.create(amount, "BBD");
}

export function fromBDT(amount: AmountInput): Money {
  return defaultManager.create(amount, "BDT");
}

export function fromBGN(amount: AmountInput): Money {
  return defaultManager.create(amount, "BGN");
}

export function fromBHD(amount: AmountInput): Money {
  return defaultManager.create(amount, "BHD");
}

export function fromBIF(amount: AmountInput): Money {
  return defaultManager.create(amount, "BIF");
}

export function fromBMD(amount: AmountInput): Money {
  return defaultManager.create(amount, "BMD");
}

export function fromBND(amount: AmountInput): Money {
  return defaultManager.create(amount, "BND");
}

export function fromBOB(amount: AmountInput): Money {
  return defaultManager.create(amount, "BOB");
}

export function fromBOV(amount: AmountInput): Money {
  return defaultManager.create(amount, "BOV");
}

export function fromBRL(amount: AmountInput): Money {
  return defaultManager.create(amount, "BRL");
}

export function fromBSD(amount: AmountInput): Money {
  return defaultManager.create(amount, "BSD");
}

export function fromBTN(amount: AmountInput): Money {
  return defaultManager.create(amount, "BTN");
}

export function fromBWP(amount: AmountInput): Money {
  return defaultManager.create(amount, "BWP");
}

export function fromBYN(amount: AmountInput): Money {
  return defaultManager.create(amount, "BYN");
}

export function fromBZD(amount: AmountInput): Money {
  return defaultManager.create(amount, "BZD");
}

export function fromCAD(amount: AmountInput): Money {
  return defaultManager.create(amount, "CAD");
}

export function fromCDF(amount: AmountInput): Money {
  return defaultManager.create(amount, "CDF");
}

export function fromCHE(amount: AmountInput): Money {
  return defaultManager.create(amount, "CHE");
}

export function fromCHF(amount: AmountInput): Money {
  return defaultManager.create(amount, "CHF");
}

export function fromCHW(amount: AmountInput): Money {
  return defaultManager.create(amount, "CHW");
}

export function fromCLF(amount: AmountInput): Money {
  return defaultManager.create(amount, "CLF");
}

export function fromCLP(amount: AmountInput): Money {
  return defaultManager.create(amount, "CLP");
}

export function fromCNY(amount: AmountInput): Money {
  return defaultManager.create(amount, "CNY");
}

export function fromCOP(amount: AmountInput): Money {
  return defaultManager.create(amount, "COP");
}

export function fromCOU(amount: AmountInput): Money {
  return defaultManager.create(amount, "COU");
}

export function fromCRC(amount: AmountInput): Money {
  return defaultManager.create(amount, "CRC");
}

export function fromCUC(amount: AmountInput): Money {
  return defaultManager.create(amount, "CUC");
}

export function fromCUP(amount: AmountInput): Money {
  return defaultManager.create(amount, "CUP");
}

export function fromCVE(amount: AmountInput): Money {
  return defaultManager.create(amount, "CVE");
}

export function fromCZK(amount: AmountInput): Money {
  return defaultManager.create(amount, "CZK");
}

export function fromDJF(amount: AmountInput): Money {
  return defaultManager.create(amount, "DJF");
}

export function fromDKK(amount: AmountInput): Money {
  return defaultManager.create(amount, "DKK");
}

export function fromDOP(amount: AmountInput): Money {
  return defaultManager.create(amount, "DOP");
}

export function fromDZD(amount: AmountInput): Money {
  return defaultManager.create(amount, "DZD");
}

export function fromEGP(amount: AmountInput): Money {
  return defaultManager.create(amount, "EGP");
}

export function fromERN(amount: AmountInput): Money {
  return defaultManager.create(amount, "ERN");
}

export function fromETB(amount: AmountInput): Money {
  return defaultManager.create(amount, "ETB");
}

export function fromEUR(amount: AmountInput): Money {
  return defaultManager.create(amount, "EUR");
}

export function fromFJD(amount: AmountInput): Money {
  return defaultManager.create(amount, "FJD");
}

export function fromFKP(amount: AmountInput): Money {
  return defaultManager.create(amount, "FKP");
}

export function fromGBP(amount: AmountInput): Money {
  return defaultManager.create(amount, "GBP");
}

export function fromGEL(amount: AmountInput): Money {
  return defaultManager.create(amount, "GEL");
}

export function fromGHS(amount: AmountInput): Money {
  return defaultManager.create(amount, "GHS");
}

export function fromGIP(amount: AmountInput): Money {
  return defaultManager.create(amount, "GIP");
}

export function fromGMD(amount: AmountInput): Money {
  return defaultManager.create(amount, "GMD");
}

export function fromGNF(amount: AmountInput): Money {
  return defaultManager.create(amount, "GNF");
}

export function fromGTQ(amount: AmountInput): Money {
  return defaultManager.create(amount, "GTQ");
}

export function fromGYD(amount: AmountInput): Money {
  return defaultManager.create(amount, "GYD");
}

export function fromHKD(amount: AmountInput): Money {
  return defaultManager.create(amount, "HKD");
}

export function fromHNL(amount: AmountInput): Money {
  return defaultManager.create(amount, "HNL");
}

export function fromHTG(amount: AmountInput): Money {
  return defaultManager.create(amount, "HTG");
}

export function fromHUF(amount: AmountInput): Money {
  return defaultManager.create(amount, "HUF");
}

export function fromIDR(amount: AmountInput): Money {
  return defaultManager.create(amount, "IDR");
}

export function fromILS(amount: AmountInput): Money {
  return defaultManager.create(amount, "ILS");
}

export function fromINR(amount: AmountInput): Money {
  return defaultManager.create(amount, "INR");
}

export function fromIQD(amount: AmountInput): Money {
  return defaultManager.create(amount, "IQD");
}

export function fromIRR(amount: AmountInput): Money {
  return defaultManager.create(amount, "IRR");
}

export function fromISK(amount: AmountInput): Money {
  return defaultManager.create(amount, "ISK");
}

export function fromJMD(amount: AmountInput): Money {
  return defaultManager.create(amount, "JMD");
}

export function fromJOD(amount: AmountInput): Money {
  return defaultManager.create(amount, "JOD");
}

export function fromJPY(amount: AmountInput): Money {
  return defaultManager.create(amount, "JPY");
}

export function fromKES(amount: AmountInput): Money {
  return defaultManager.create(amount, "KES");
}

export function fromKGS(amount: AmountInput): Money {
  return defaultManager.create(amount, "KGS");
}

export function fromKHR(amount: AmountInput): Money {
  return defaultManager.create(amount, "KHR");
}

export function fromKMF(amount: AmountInput): Money {
  return defaultManager.create(amount, "KMF");
}

export function fromKPW(amount: AmountInput): Money {
  return defaultManager.create(amount, "KPW");
}

export function fromKRW(amount: AmountInput): Money {
  return defaultManager.create(amount, "KRW");
}

export function fromKWD(amount: AmountInput): Money {
  return defaultManager.create(amount, "KWD");
}

export function fromKYD(amount: AmountInput): Money {
  return defaultManager.create(amount, "KYD");
}

export function fromKZT(amount: AmountInput): Money {
  return defaultManager.create(amount, "KZT");
}

export function fromLAK(amount: AmountInput): Money {
  return defaultManager.create(amount, "LAK");
}

export function fromLBP(amount: AmountInput): Money {
  return defaultManager.create(amount, "LBP");
}

export function fromLKR(amount: AmountInput): Money {
  return defaultManager.create(amount, "LKR");
}

export function fromLRD(amount: AmountInput): Money {
  return defaultManager.create(amount, "LRD");
}

export function fromLSL(amount: AmountInput): Money {
  return defaultManager.create(amount, "LSL");
}

export function fromLYD(amount: AmountInput): Money {
  return defaultManager.create(amount, "LYD");
}

export function fromMAD(amount: AmountInput): Money {
  return defaultManager.create(amount, "MAD");
}

export function fromMDL(amount: AmountInput): Money {
  return defaultManager.create(amount, "MDL");
}

export function fromMGA(amount: AmountInput): Money {
  return defaultManager.create(amount, "MGA");
}

export function fromMKD(amount: AmountInput): Money {
  return defaultManager.create(amount, "MKD");
}

export function fromMMK(amount: AmountInput): Money {
  return defaultManager.create(amount, "MMK");
}

export function fromMNT(amount: AmountInput): Money {
  return defaultManager.create(amount, "MNT");
}

export function fromMOP(amount: AmountInput): Money {
  return defaultManager.create(amount, "MOP");
}

export function fromMUR(amount: AmountInput): Money {
  return defaultManager.create(amount, "MUR");
}

export function fromMRU(amount: AmountInput): Money {
  return defaultManager.create(amount, "MRU");
}

export function fromMVR(amount: AmountInput): Money {
  return defaultManager.create(amount, "MVR");
}

export function fromMWK(amount: AmountInput): Money {
  return defaultManager.create(amount, "MWK");
}

export function fromMXN(amount: AmountInput): Money {
  return defaultManager.create(amount, "MXN");
}

export function fromMXV(amount: AmountInput): Money {
  return defaultManager.create(amount, "MXV");
}

export function fromMYR(amount: AmountInput): Money {
  return defaultManager.create(amount, "MYR");
}

export function fromMZN(amount: AmountInput): Money {
  return defaultManager.create(amount, "MZN");
}

export function fromNAD(amount: AmountInput): Money {
  return defaultManager.create(amount, "NAD");
}

export function fromNGN(amount: AmountInput): Money {
  return defaultManager.create(amount, "NGN");
}

export function fromNIO(amount: AmountInput): Money {
  return defaultManager.create(amount, "NIO");
}

export function fromNOK(amount: AmountInput): Money {
  return defaultManager.create(amount, "NOK");
}

export function fromNPR(amount: AmountInput): Money {
  return defaultManager.create(amount, "NPR");
}

export function fromNZD(amount: AmountInput): Money {
  return defaultManager.create(amount, "NZD");
}

export function fromOMR(amount: AmountInput): Money {
  return defaultManager.create(amount, "OMR");
}

export function fromPAB(amount: AmountInput): Money {
  return defaultManager.create(amount, "PAB");
}

export function fromPEN(amount: AmountInput): Money {
  return defaultManager.create(amount, "PEN");
}

export function fromPGK(amount: AmountInput): Money {
  return defaultManager.create(amount, "PGK");
}

export function fromPHP(amount: AmountInput): Money {
  return defaultManager.create(amount, "PHP");
}

export function fromPKR(amount: AmountInput): Money {
  return defaultManager.create(amount, "PKR");
}

export function fromPLN(amount: AmountInput): Money {
  return defaultManager.create(amount, "PLN");
}

export function fromPYG(amount: AmountInput): Money {
  return defaultManager.create(amount, "PYG");
}

export function fromQAR(amount: AmountInput): Money {
  return defaultManager.create(amount, "QAR");
}

export function fromRON(amount: AmountInput): Money {
  return defaultManager.create(amount, "RON");
}

export function fromRSD(amount: AmountInput): Money {
  return defaultManager.create(amount, "RSD");
}

export function fromRUB(amount: AmountInput): Money {
  return defaultManager.create(amount, "RUB");
}

export function fromRWF(amount: AmountInput): Money {
  return defaultManager.create(amount, "RWF");
}

export function fromSAR(amount: AmountInput): Money {
  return defaultManager.create(amount, "SAR");
}

export function fromSBD(amount: AmountInput): Money {
  return defaultManager.create(amount, "SBD");
}

export function fromSCR(amount: AmountInput): Money {
  return defaultManager.create(amount, "SCR");
}

export function fromSDG(amount: AmountInput): Money {
  return defaultManager.create(amount, "SDG");
}

export function fromSEK(amount: AmountInput): Money {
  return defaultManager.create(amount, "SEK");
}

export function fromSGD(amount: AmountInput): Money {
  return defaultManager.create(amount, "SGD");
}

export function fromSHP(amount: AmountInput): Money {
  return defaultManager.create(amount, "SHP");
}

export function fromSLE(amount: AmountInput): Money {
  return defaultManager.create(amount, "SLE");
}

export function fromSLL(amount: AmountInput): Money {
  return defaultManager.create(amount, "SLL");
}

export function fromSOS(amount: AmountInput): Money {
  return defaultManager.create(amount, "SOS");
}

export function fromSRD(amount: AmountInput): Money {
  return defaultManager.create(amount, "SRD");
}

export function fromSSP(amount: AmountInput): Money {
  return defaultManager.create(amount, "SSP");
}

export function fromSTN(amount: AmountInput): Money {
  return defaultManager.create(amount, "STN");
}

export function fromSVC(amount: AmountInput): Money {
  return defaultManager.create(amount, "SVC");
}

export function fromSYP(amount: AmountInput): Money {
  return defaultManager.create(amount, "SYP");
}

export function fromSZL(amount: AmountInput): Money {
  return defaultManager.create(amount, "SZL");
}

export function fromTHB(amount: AmountInput): Money {
  return defaultManager.create(amount, "THB");
}

export function fromTJS(amount: AmountInput): Money {
  return defaultManager.create(amount, "TJS");
}

export function fromTMT(amount: AmountInput): Money {
  return defaultManager.create(amount, "TMT");
}

export function fromTND(amount: AmountInput): Money {
  return defaultManager.create(amount, "TND");
}

export function fromTOP(amount: AmountInput): Money {
  return defaultManager.create(amount, "TOP");
}

export function fromTRY(amount: AmountInput): Money {
  return defaultManager.create(amount, "TRY");
}

export function fromTTD(amount: AmountInput): Money {
  return defaultManager.create(amount, "TTD");
}

export function fromTWD(amount: AmountInput): Money {
  return defaultManager.create(amount, "TWD");
}

export function fromTZS(amount: AmountInput): Money {
  return defaultManager.create(amount, "TZS");
}

export function fromUAH(amount: AmountInput): Money {
  return defaultManager.create(amount, "UAH");
}

export function fromUGX(amount: AmountInput): Money {
  return defaultManager.create(amount, "UGX");
}

export function fromUSD(amount: AmountInput): Money {
  return defaultManager.create(amount, "USD");
}

export function fromUSN(amount: AmountInput): Money {
  return defaultManager.create(amount, "USN");
}

export function fromUYI(amount: AmountInput): Money {
  return defaultManager.create(amount, "UYI");
}

export function fromUYU(amount: AmountInput): Money {
  return defaultManager.create(amount, "UYU");
}

export function fromUYW(amount: AmountInput): Money {
  return defaultManager.create(amount, "UYW");
}

export function fromUZS(amount: AmountInput): Money {
  return defaultManager.create(amount, "UZS");
}

export function fromVES(amount: AmountInput): Money {
  return defaultManager.create(amount, "VES");
}

export function fromVND(amount: AmountInput): Money {
  return defaultManager.create(amount, "VND");
}

export function fromVUV(amount: AmountInput): Money {
  return defaultManager.create(amount, "VUV");
}

export function fromWST(amount: AmountInput): Money {
  return defaultManager.create(amount, "WST");
}

export function fromXAF(amount: AmountInput): Money {
  return defaultManager.create(amount, "XAF");
}

export function fromXAG(amount: AmountInput): Money {
  return defaultManager.create(amount, "XAG");
}

export function fromXAU(amount: AmountInput): Money {
  return defaultManager.create(amount, "XAU");
}

export function fromXBA(amount: AmountInput): Money {
  return defaultManager.create(amount, "XBA");
}

export function fromXBB(amount: AmountInput): Money {
  return defaultManager.create(amount, "XBB");
}

export function fromXBC(amount: AmountInput): Money {
  return defaultManager.create(amount, "XBC");
}

export function fromXBD(amount: AmountInput): Money {
  return defaultManager.create(amount, "XBD");
}

export function fromXCD(amount: AmountInput): Money {
  return defaultManager.create(amount, "XCD");
}

export function fromXDR(amount: AmountInput): Money {
  return defaultManager.create(amount, "XDR");
}

export function fromXOF(amount: AmountInput): Money {
  return defaultManager.create(amount, "XOF");
}

export function fromXPD(amount: AmountInput): Money {
  return defaultManager.create(amount, "XPD");
}

export function fromXPF(amount: AmountInput): Money {
  return defaultManager.create(amount, "XPF");
}

export function fromXPT(amount: AmountInput): Money {
  return defaultManager.create(amount, "XPT");
}

export function fromXSU(amount: AmountInput): Money {
  return defaultManager.create(amount, "XSU");
}

export function fromXTS(amount: AmountInput): Money {
  return defaultManager.create(amount, "XTS");
}

export function fromXUA(amount: AmountInput): Money {
  return defaultManager.create(amount, "XUA");
}

export function fromXXX(amount: AmountInput): Money {
  return defaultManager.create(amount, "XXX");
}

export function fromYER(amount: AmountInput): Money {
  return defaultManager.create(amount, "YER");
}

export function fromZAR(amount: AmountInput): Money {
  return defaultManager.create(amount, "ZAR");
}

export function fromZMW(amount: AmountInput): Money {
  return defaultManager.create(amount, "ZMW");
}

export function fromZWL(amount: AmountInput): Money {
  return defaultManager.create(amount, "ZWL");
}
