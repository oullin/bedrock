# lottery

Probabilistic execution helpers split out from the old `support` module.

## Overview

The `lottery` module exposes the `Lottery` type and its helpers for
probabilistic execution and deterministic test control.

**Module:** `github.com/bedrock/packages/lottery`

```bash
go get github.com/bedrock/packages/lottery@latest
```

## Usage

```go
lottery.NewLottery(1, 100).Winner(func(...any) any {
    return nil
}).Loser(func(...any) any {
    return nil
}).Run()

lottery.LotteryOdds(1, 10).Choose()
```
