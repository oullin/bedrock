package spark_test

import (
	"net/http"
	"testing"

	"github.com/bedrock/packages/spark"
)

type stubBillable struct {
	id    int64
	btype string
}

func (s stubBillable) BillableID() int64     { return s.id }
func (s stubBillable) BillableType() string  { return s.btype }
func (s stubBillable) BillableUUID() string  { return "uuid-123" }
func (s stubBillable) BillableName() string  { return "Test Team" }
func (s stubBillable) BillableEmail() string { return "team@example.com" }

func TestManagerBillableRegistration(t *testing.T) {
	m := spark.NewManager()

	m.Billable("team").
		Resolve(func(r *http.Request) (spark.Billable, error) {
			return stubBillable{id: 1, btype: "team"}, nil
		}).
		Authorize(func(b spark.Billable, r *http.Request) bool {
			return b.BillableID() == 1
		})

	b, err := m.ResolveBillable("team", nil)

	if err != nil {
		t.Fatalf("ResolveBillable() error = %v", err)
	}

	if b.BillableID() != 1 {
		t.Errorf("BillableID() = %d, want 1", b.BillableID())
	}

	if !m.IsAuthorized(b, nil) {
		t.Error("IsAuthorized() = false, want true")
	}
}

func TestManagerPlans(t *testing.T) {
	m := spark.NewManager()

	m.AddPlan("team", spark.SparkPlan{ID: "pri_monthly", Name: "Pro", Interval: "monthly"})
	m.AddPlan("team", spark.SparkPlan{ID: "pri_yearly", Name: "Pro", Interval: "yearly"})

	plans := m.Plans("team")

	if len(plans) != 2 {
		t.Fatalf("Plans() returned %d plans, want 2", len(plans))
	}
}

func TestManagerDefaultBillableType(t *testing.T) {
	m := spark.NewManager()
	m.RegisterBillable(spark.BillableConfig{ModelName: "team"})

	if got := m.DefaultBillableType(); got != "team" {
		t.Errorf("DefaultBillableType() = %q, want %q", got, "team")
	}

	m.RegisterBillable(spark.BillableConfig{ModelName: "user"})

	if got := m.DefaultBillableType(); got != "user" {
		t.Errorf("DefaultBillableType() with two = %q, want %q", got, "user")
	}
}

func TestManagerPerSeatBilling(t *testing.T) {
	m := spark.NewManager()

	m.Billable("team").ChargePerSeat("Team Member", func(b spark.Billable) int { return 5 })

	if !m.ChargesPerSeat("team") {
		t.Error("ChargesPerSeat() = false, want true")
	}

	if got := m.SeatCount("team", stubBillable{id: 1, btype: "team"}); got != 5 {
		t.Errorf("SeatCount() = %d, want 5", got)
	}

	if got := m.SeatName("team"); got != "Team Member" {
		t.Errorf("SeatName() = %q, want %q", got, "Team Member")
	}
}
