package handler_test

import (
	"context"

	"github.com/bedrock/packages/billing"
)

type stubBillable struct {
	id       int64
	btype    string
	name     string
	email    string
	provider string
}

type testSubStore struct {
	subs []*billing.Subscription
}

type testOrderStore struct{}

type testProductStore struct {
	subProducts     []billing.Product
	oneTimeProducts []billing.Product
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

func (s *testSubStore) FindByID(context.Context, int64) (*billing.Subscription, error) {
	return nil, nil
}
func (s *testSubStore) FindByProviderID(context.Context, string) (*billing.Subscription, error) {
	return nil, nil
}
func (s *testSubStore) CurrentForBillable(context.Context, string, int64) (*billing.Subscription, error) {
	return nil, nil
}
func (s *testSubStore) ActiveForBillable(context.Context, string, int64) ([]*billing.Subscription, error) {
	return s.subs, nil
}
func (s *testSubStore) Create(context.Context, *billing.Subscription) error { return nil }
func (s *testSubStore) Save(context.Context, *billing.Subscription) error   { return nil }
func (s *testSubStore) Delete(context.Context, int64) error                 { return nil }

func (s *testOrderStore) FindByID(context.Context, int64) (*billing.Order, error) {
	return nil, nil
}
func (s *testOrderStore) FindByBillable(context.Context, int64, int) ([]billing.Order, error) {
	return nil, nil
}
func (s *testOrderStore) Create(context.Context, *billing.Order) error { return nil }
func (s *testOrderStore) Save(context.Context, *billing.Order) error   { return nil }
func (s *testOrderStore) HasCompletedForProduct(context.Context, int64, int64) (bool, error) {
	return false, nil
}

func (s *testProductStore) FindByID(context.Context, int64) (*billing.Product, error) {
	return nil, nil
}
func (s *testProductStore) Active(context.Context) ([]billing.Product, error) {
	return append(s.subProducts, s.oneTimeProducts...), nil
}
func (s *testProductStore) ActiveSubscriptions(context.Context) ([]billing.Product, error) {
	return s.subProducts, nil
}
func (s *testProductStore) ActiveOneTime(context.Context) ([]billing.Product, error) {
	return s.oneTimeProducts, nil
}
