# money

<!-- ref: @bedrock/code-0180 -->
<!-- ref: @bedrock/code-0093 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

The money package provides Bedrock's Go implementation for this surface.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/money@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/money/...
```

## Source Coverage

| Package      | Purpose                                        |
| ------------ | ---------------------------------------------- |
| `calculator` | Public calculator API surface for this module. |
| `currency`   | Public currency API surface for this module.   |
| `exception`  | Public exception API surface for this module.  |
| `exchange`   | Public exchange API surface for this module.   |
| `format`     | Public format API surface for this module.     |
| `money`      | Public money API surface for this module.      |
| `parser`     | Public parser API surface for this module.     |

## Core Concepts

The money reference is organized around the exported Go surface for package `money`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Aggregator`, `Amount`, `Calculator`, `Converter`, `Currency`, `DefaultProvider`, `Exchange`, `Formatter`, `ISOCodePattern`, `JSON`, `JSONRawData`, `Manager`, `Map`, `Money`, `Parser`, `Provider`, `Symbol`                                                                                                                                                                                                                                                                                                                             |
| Constructors and functions | `Absolute`, `Add`, `AddFrom`, `AddRate`, `All`, `Allocate`, `Amount`, `AsMajorUnits`, `AssertSameCurrency`, `Avg`, `Compare`, `CompareAmount`, `Convert`, `ConvertAmount`, `ConvertAmountWithRate`, `ConvertWithRate`, `Create`, `CreateFromFloat`, `CreateFromString`, `Currency`, and 260 more                                                                                                                                                                                                                                          |
| Variables                  | `ErrCurrencyConversionNotFound`, `ErrCurrencyMismatch`, `ErrCurrencyNotFound`, `ErrCurrencyNotSpecified`, `ErrEmptyAmountString`, `ErrInvalidAggregatorProvider`, `ErrInvalidAmount`, `ErrInvalidAmountFraction`, `ErrInvalidAmountMultiple`, `ErrInvalidExchangeRate`, `ErrInvalidJSONUnmarshal`, `ErrInvalidMoneyString`, `ErrInvalidSplit`, `ErrJSONMarshalFuncNil`, `ErrJSONUnmarshalFuncNil`, `ErrNegativeRatios`, `ErrNoConverterProvided`, `ErrNoCurrencyInstance`, `ErrNoCurrencyManager`, `ErrNoCurrencyMapDataset`, and 10 more |
| Constants                  | `AED`, `AFN`, `ALL`, `AMD`, `ANG`, `AOA`, `ARS`, `AUD`, `AWG`, `AZN`, `BAM`, `BBD`, `BDT`, `BGN`, `BHD`, `BIF`, `BMD`, `BND`, `BOB`, `BOV`, and 160 more                                                                                                                                                                                                                                                                                                                                                                                  |

### Capability Matrix

| Capability                         | Documentation note                                                                                                   |
| ---------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers               | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/money"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/money` cover the supported creation paths, default values, and parity behavior.

## Configuration

Bedrock documents behavior through Go options and constructor arguments:

| Upstream shape    | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/money/...
```

## API Reference

### Exported Types

| Type              | Notes                                                                              |
| ----------------- | ---------------------------------------------------------------------------------- |
| `Aggregator`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Amount`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Calculator`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Converter`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Currency`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exchange`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Formatter`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ISOCodePattern`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSON`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSONRawData`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Map`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Money`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parser`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provider`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Symbol`          | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                         | Notes                                                                              |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `Absolute`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Add`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddFrom`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddRate`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Allocate`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Amount`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AsMajorUnits`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSameCurrency`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Avg`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Compare`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompareAmount`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Convert`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConvertAmount`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConvertAmountWithRate`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConvertWithRate`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Create`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateFromFloat`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateFromString`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Currency`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DbScan`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DbValue`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Display`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Divide`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Equals`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindByCode`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindByNumericCode`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Format`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Formatter`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromAED`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromAFN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromALL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromAMD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromANG`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromAOA`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromARS`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromAUD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromAWG`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromAZN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBAM`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBBD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBDT`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBGN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBHD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBIF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBMD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBND`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBOB`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBOV`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBRL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBSD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBTN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBWP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBYN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBZD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCAD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCDF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCHE`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCHF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCHW`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCLF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCLP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCNY`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCOP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCOU`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCRC`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCUC`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCUP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCVE`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromCZK`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromDJF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromDKK`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromDOP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromDZD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromEGP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromERN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromETB`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromEUR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromFJD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromFKP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromGBP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromGEL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromGHS`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromGIP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromGMD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromGNF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromGTQ`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromGYD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromHKD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromHNL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromHTG`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromHUF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromIDR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromILS`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromINR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromIQD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromIRR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromISK`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromJMD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromJOD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromJPY`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromKES`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromKGS`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromKHR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromKMF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromKPW`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromKRW`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromKWD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromKYD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromKZT`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromLAK`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromLBP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromLKR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromLRD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromLSL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromLYD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMAD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMDL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMGA`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMKD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMMK`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMNT`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMOP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMRU`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMUR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMVR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMWK`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMXN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMXV`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMYR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromMZN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromNAD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromNGN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromNIO`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromNOK`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromNPR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromNZD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromOMR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromPAB`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromPEN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromPGK`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromPHP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromPKR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromPLN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromPYG`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromQAR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromRON`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromRSD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromRUB`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromRWF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSAR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSBD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSCR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSDG`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSEK`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSGD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSHP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSLE`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSLL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSOS`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSRD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSSP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSTN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSVC`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSYP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSZL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromTHB`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromTJS`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromTMT`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromTND`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromTOP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromTRY`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromTTD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromTWD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromTZS`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromUAH`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromUGX`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromUSD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromUSN`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromUYI`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromUYU`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromUYW`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromUZS`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromVES`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromVND`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromVUV`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromWST`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXAF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXAG`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXAU`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXBA`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXBB`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXBC`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXBD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXCD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXDR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXOF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXPD`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXPF`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXPT`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXSU`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXTS`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXUA`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromXXX`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromYER`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromZAR`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromZMW`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromZWL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCode`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCodes`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCurrencyManager`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDBMoneyValueSeparator`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefault`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetExchange`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPattern`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRate`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetSymbols`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetSymbolsLongestFirst`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GreaterThan`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GreaterThanOrEqual`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasInvalidState`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmpty`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsInvalid`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNegative`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNotEmpty`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPositive`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsValid`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsZero`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LessThan`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LessThanOrEqual`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Marshal`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarshalJSON`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Max`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Min`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Modulus`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Multiply`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Negative`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAggregator`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCalculator`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConverter`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewErrInvalidJSONUnmarshalFrom` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewExchange`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFormatter`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewISOCodePattern`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewISOCodePatternWith`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewJson`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewJsonWithParser`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManagerFor`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManagerWith`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewParser`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewParserWith`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseAmount`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseAmountString`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseAmountWithDecimalComma`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseDecimal`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseDecimalParts`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseDecimalWithComma`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseStringSign`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Ration`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolve`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Round`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SafeAdd`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SafeMultiply`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SafeSubtract`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SameCurrency`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scan`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCurrency`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDBMoneyValueSeparator`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetMarshal`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetUnmarshal`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Split`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subtract`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sum`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMajorUnits`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unmarshal`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnmarshalJSON`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidateAndPadDecimal`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Value`                          | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                            | Notes                                                                              |
| ------------------------------- | ---------------------------------------------------------------------------------- |
| `AED`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AFN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ALL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AMD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ANG`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AOA`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ARS`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AUD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AWG`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AZN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BAM`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BBD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BDT`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BGN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BHD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BIF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BMD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BND`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BOB`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BOV`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BRL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BSD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BTN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BWP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BYN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BZD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CAD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CDF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CHE`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CHF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CHW`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CLF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CLP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CNY`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `COP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `COU`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CRC`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CUC`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CUP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CVE`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CZK`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DJF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DKK`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DOP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DZD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultDBMoneyValueSeparator`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EGP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ERN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ETB`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EUR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrCurrencyConversionNotFound` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrCurrencyMismatch`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrCurrencyNotFound`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrCurrencyNotSpecified`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrEmptyAmountString`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidAggregatorProvider`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidAmount`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidAmountFraction`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidAmountMultiple`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidExchangeRate`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidJSONUnmarshal`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidMoneyString`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidSplit`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrJSONMarshalFuncNil`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrJSONUnmarshalFuncNil`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNegativeRatios`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoConverterProvided`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoCurrencyInstance`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoCurrencyManager`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoCurrencyMapDataset`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoJSONParserProvided`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoMoneyProvided`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoMultipliersProvided`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoRatiosProvided`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrOverflow`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrParserInvalidState`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrParserNotProvided`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrRatiosExceedMaxInt`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FJD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FKP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GBP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GEL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GHS`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GIP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GMD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GNF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GTQ`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GYD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HKD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HNL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HTG`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HUF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IDR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ILS`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `INR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IQD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IRR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ISK`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JMD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JOD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JPY`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KES`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KGS`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KHR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KMF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KPW`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KRW`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KWD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KYD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KZT`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LAK`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LBP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LKR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LRD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LSL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LYD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MAD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MDL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MGA`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MKD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MMK`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MNT`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MOP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MRU`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MUR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MVR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MWK`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MXN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MXV`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MYR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MZN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NAD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NGN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NIO`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NOK`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NPR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NZD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCurrenciesMap`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCurrenciesMapFrom`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OMR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PAB`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PEN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PGK`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PHP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PKR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PLN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PYG`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QAR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RON`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RSD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RUB`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RWF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SAR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SBD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SCR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SDG`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SEK`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SGD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SHP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SLE`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SLL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SOS`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SRD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SSP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `STN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SVC`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SYP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SZL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `THB`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TJS`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TMT`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TND`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TOP`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TRY`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TTD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TWD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TZS`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UAH`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UGX`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `USD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `USN`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UYI`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UYU`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UYW`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UZS`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VES`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VND`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VUV`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WST`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XAF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XAG`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XAU`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XBA`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XBB`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XBC`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XBD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XCD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XDR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XOF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XPD`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XPF`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XPT`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XSU`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XTS`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XUA`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `XXX`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `YER`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZAR`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZMW`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZWL`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
