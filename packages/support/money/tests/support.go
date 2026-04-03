package tests

import (
	"math"
	"math/big"
	"testing"
)

func TestRequire[T any](t testing.TB, fn func() (T, error)) T {
	t.Helper()

	value, err := fn()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return value
}

func TestRequireNoErr(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPow10i64Money(t testing.TB, exp int) int64 {
	t.Helper()

	if exp < 0 || exp > 18 {
		t.Fatalf("exponent %d out of safe int64 range [0, 18]", exp)
	}

	n := int64(1)

	for i := 0; i < exp; i++ {
		n *= 10
	}

	return n
}

func TestAbs64Money(t testing.TB, v int64) int64 {
	t.Helper()

	if v == math.MinInt64 {
		t.Fatalf("cannot compute absolute value of math.MinInt64")
	}

	if v < 0 {
		return -v
	}

	return v
}

func TestMustParseFloatMoney(t testing.TB, s string) float64 {
	t.Helper()

	f, ok := new(big.Rat).SetString(s)

	if !ok {
		t.Fatalf("invalid decimal rate string %q", s)
	}

	f64, _ := f.Float64()

	return f64
}

func TestExpectedConvertAmountMoney(t testing.TB, amount int64, fromFraction, toFraction int, rate string) int64 {
	t.Helper()

	rateRat, ok := new(big.Rat).SetString(rate)

	if !ok {
		t.Fatalf("invalid rate: %q", rate)

		return 0
	}

	scaled := new(big.Rat).SetInt64(amount)

	scaled.Mul(scaled, rateRat)
	scaled.Mul(scaled, big.NewRat(TestPow10i64Money(t, toFraction), 1))
	scaled.Quo(scaled, big.NewRat(TestPow10i64Money(t, fromFraction), 1))

	return TestRoundRatHalfAwayFromZeroMoney(t, scaled)
}

func TestRoundRatHalfAwayFromZeroMoney(t testing.TB, r *big.Rat) int64 {
	t.Helper()

	if r.Sign() == 0 {
		return 0
	}

	num := new(big.Int).Set(r.Num())
	den := new(big.Int).Set(r.Denom())

	neg := num.Sign() < 0

	if neg {
		num.Abs(num)
	}

	q, rem := new(big.Int), new(big.Int)
	q.QuoRem(num, den, rem)

	twiceRem := new(big.Int).Lsh(rem, 1)

	if twiceRem.Cmp(den) >= 0 {
		q.Add(q, big.NewInt(1))
	}

	if neg {
		q.Neg(q)
	}

	return q.Int64()
}
