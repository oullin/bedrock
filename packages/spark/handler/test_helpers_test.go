package handler_test

import (
	"context"

	"github.com/bedrock/packages/spark"
)

type stubBillable struct {
	id       int64
	btype    string
	name     string
	email    string
	provider string
}

type testSubStore struct {
	subs []*spark.Subscription
}

type testOrderStore struct{}

type testProductStore struct {
	subProducts     []spark.Product
	oneTimeProducts []spark.Product
}

func (b *stubBillable) BillableID() int64       { return b.id }
func (b *stubBillable) BillableType() string    { return b.btype }
func (b *stubBillable) BillableName() string    { return b.name }
func (b *stubBillable) BillableEmail() string   { return b.email }
func (b *stubBillable) PaymentProvider() string { return b.provider }
func (b *stubBillable) SetPaymentProvider(p string) error {
	b.provider = p

	return nil
}

func (s *testSubStore) FindByID(context.Context, int64) (*spark.Subscription, error) {
	return nil, nil
}
func (s *testSubStore) FindByProviderID(context.Context, string) (*spark.Subscription, error) {
	return nil, nil
}
func (s *testSubStore) CurrentForBillable(context.Context, string, int64) (*spark.Subscription, error) {
	return nil, nil
}
func (s *testSubStore) ActiveForBillable(context.Context, string, int64) ([]*spark.Subscription, error) {
	return s.subs, nil
}
func (s *testSubStore) Create(context.Context, *spark.Subscription) error { return nil }
func (s *testSubStore) Save(context.Context, *spark.Subscription) error   { return nil }
func (s *testSubStore) Delete(context.Context, int64) error               { return nil }

func (s *testOrderStore) FindByID(context.Context, int64) (*spark.Order, error) {
	return nil, nil
}
func (s *testOrderStore) FindByBillable(context.Context, int64, int) ([]spark.Order, error) {
	return nil, nil
}
func (s *testOrderStore) Create(context.Context, *spark.Order) error { return nil }
func (s *testOrderStore) Save(context.Context, *spark.Order) error   { return nil }
func (s *testOrderStore) HasCompletedForProduct(context.Context, int64, int64) (bool, error) {
	return false, nil
}

func (s *testProductStore) FindByID(context.Context, int64) (*spark.Product, error) {
	return nil, nil
}
func (s *testProductStore) Active(context.Context) ([]spark.Product, error) {
	return append(s.subProducts, s.oneTimeProducts...), nil
}
func (s *testProductStore) ActiveSubscriptions(context.Context) ([]spark.Product, error) {
	return s.subProducts, nil
}
func (s *testProductStore) ActiveOneTime(context.Context) ([]spark.Product, error) {
	return s.oneTimeProducts, nil
}
