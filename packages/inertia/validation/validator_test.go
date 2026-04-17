package validation

import (
	"testing"
)

type testForm struct {
	Name  string `json:"name" validate:"required,max=255"`
	Email string `json:"email" validate:"required,email,max=255"`
	Phone string `json:"phone" validate:"omitempty,max=255"`
	Age   int    `json:"age" validate:"omitempty,gte=0,lte=150"`
}

//nolint:goconst

// Should use json tag "name" not struct field "Name"

type formTagForm struct {
	Title string `form:"title" validate:"required"`
}

type minForm struct {
	Password string `json:"password" validate:"required,min=8"`
}

type lenForm struct {
	Code string `json:"code" validate:"required,len=6"`
}

type rangeForm struct {
	Score int `json:"score" validate:"gte=1,lte=100"`
}

type oneofForm struct {
	Role string `json:"role" validate:"required,oneof=admin editor viewer"`
}

type urlForm struct {
	Website string `json:"website" validate:"required,url"`
}

type eqfieldForm struct {
	Password string `json:"password" validate:"required"`
	Confirm  string `json:"confirm" validate:"required,eqfield=Password"`
}

type customTagForm struct {
	Value string `json:"value" validate:"required,alphanum"`
}

type jsonDashForm struct {
	Internal string `json:"-" validate:"required"`
}

type jsonOmitemptyForm struct {
	Name string `json:"name,omitempty" validate:"required"`
}

func TestValidate_PassesWithValidData(t *testing.T) {
	t.Parallel()

	form := testForm{
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "555-1234",
		Age:   30,
	}

	errors := Validate(form)

	if errors != nil {
		t.Errorf("expected no errors, got %v", errors)
	}
}

func TestValidate_RequiredFields(t *testing.T) {
	t.Parallel()

	form := testForm{}

	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	if _, ok := errors["name"]; !ok {
		t.Error("expected error for 'name'")
	}

	if _, ok := errors["email"]; !ok {
		t.Error("expected error for 'email'")
	}

	if _, ok := errors["phone"]; ok {
		t.Error("did not expect error for optional 'phone'")
	}
}

func TestValidate_InvalidEmail(t *testing.T) {
	t.Parallel()

	form := testForm{
		Name:  "John",
		Email: "not-an-email",
	}

	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["email"].(string)

	if !ok {
		t.Fatalf("expected string error for email, got %T", errors["email"])
	}

	if msg != "The email field must be a valid email address." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_MaxLength(t *testing.T) {
	t.Parallel()

	long := make([]byte, 256)

	for i := range long {
		long[i] = 'a'
	}

	form := testForm{
		Name:  string(long),
		Email: "valid@example.com",
	}

	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	if _, ok := errors["name"]; !ok {
		t.Error("expected max-length error for 'name'")
	}
}

func TestValidate_UsesJSONTagForFieldName(t *testing.T) {
	t.Parallel()

	form := testForm{}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors")
	}

	if _, ok := errors["name"]; !ok {
		t.Error("expected field key to be json tag 'name'")
	}

	if _, ok := errors["Name"]; ok {
		t.Error("field key should not be struct field 'Name'")
	}
}

func TestValidate_FallsBackToFormTag(t *testing.T) {
	t.Parallel()

	form := formTagForm{}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors")
	}

	if _, ok := errors["title"]; !ok {
		t.Error("expected field key to be form tag 'title'")
	}
}

func TestValidate_MaxLength_AtBoundary(t *testing.T) {
	t.Parallel()

	boundary := make([]byte, 255)

	for i := range boundary {
		boundary[i] = 'a'
	}

	form := testForm{
		Name:  string(boundary),
		Email: "valid@example.com",
	}

	errors := Validate(form)

	if errors != nil {
		t.Errorf("expected no errors for 255-char name, got %v", errors)
	}
}

func TestValidate_MaxLength_Phone(t *testing.T) {
	t.Parallel()

	long := make([]byte, 256)

	for i := range long {
		long[i] = '5'
	}

	form := testForm{
		Name:  "John",
		Email: "john@example.com",
		Phone: string(long),
	}

	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors for 256-char phone, got nil")
	}

	msg, ok := errors["phone"].(string)

	if !ok {
		t.Fatalf("expected string error for phone, got %T", errors["phone"])
	}

	if msg != "The phone field must not exceed 255 characters." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestHumanize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"first_name", "first name"},
		{"firstName", "first name"},
		{"FirstName", "first name"},
		{"email", "email"},
		{"already lower", "already lower"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := humanize(tt.input)

			if got != tt.want {
				t.Errorf("humanize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidate_ReturnsNilForValidStruct(t *testing.T) {
	t.Parallel()

	form := testForm{
		Name:  "Valid",
		Email: "valid@example.com",
	}

	errors := Validate(form)

	if errors != nil {
		t.Errorf("expected nil, got %v", errors)
	}
}

func TestValidate_MinMessage(t *testing.T) {
	t.Parallel()

	form := minForm{Password: "short"}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["password"].(string)

	if !ok {
		t.Fatalf("expected string error, got %T", errors["password"])
	}

	if msg != "The password field must be at least 8 characters." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_LenMessage(t *testing.T) {
	t.Parallel()

	form := lenForm{Code: "abc"}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["code"].(string)

	if !ok {
		t.Fatalf("expected string error, got %T", errors["code"])
	}

	if msg != "The code field must be exactly 6 characters." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_GteMessage(t *testing.T) {
	t.Parallel()

	form := rangeForm{Score: 0}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["score"].(string)

	if !ok {
		t.Fatalf("expected string error, got %T", errors["score"])
	}

	if msg != "The score field must be at least 1." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_LteMessage(t *testing.T) {
	t.Parallel()

	form := rangeForm{Score: 101}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["score"].(string)

	if !ok {
		t.Fatalf("expected string error, got %T", errors["score"])
	}

	if msg != "The score field must be at most 100." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_OneofMessage(t *testing.T) {
	t.Parallel()

	form := oneofForm{Role: "superuser"}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["role"].(string)

	if !ok {
		t.Fatalf("expected string error, got %T", errors["role"])
	}

	if msg != "The role field must be one of: admin editor viewer." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_URLMessage(t *testing.T) {
	t.Parallel()

	form := urlForm{Website: "not-a-url"}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["website"].(string)

	if !ok {
		t.Fatalf("expected string error, got %T", errors["website"])
	}

	if msg != "The website field must be a valid URL." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_EqfieldMessage(t *testing.T) {
	t.Parallel()

	form := eqfieldForm{Password: "secret123", Confirm: "different"}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["confirm"].(string)

	if !ok {
		t.Fatalf("expected string error, got %T", errors["confirm"])
	}

	if msg != "The confirm field must match Password." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_DefaultMessage(t *testing.T) {
	t.Parallel()

	form := customTagForm{Value: "not alpha-num!@#"}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	msg, ok := errors["value"].(string)

	if !ok {
		t.Fatalf("expected string error, got %T", errors["value"])
	}

	if msg != "The value field is invalid." {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestValidate_JSONDashFallsBackToStructName(t *testing.T) {
	t.Parallel()

	form := jsonDashForm{}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	if _, ok := errors["Internal"]; !ok {
		t.Error("expected field key to be struct field name 'Internal' when json tag is '-'")
	}
}

func TestValidate_JSONOmitemptyStripsComma(t *testing.T) {
	t.Parallel()

	form := jsonOmitemptyForm{}
	errors := Validate(form)

	if errors == nil {
		t.Fatal("expected errors, got nil")
	}

	if _, ok := errors["name"]; !ok {
		t.Error("expected field key 'name', stripping ',omitempty'")
	}
}

func TestValidate_NonStructInput(t *testing.T) {
	t.Parallel()

	errors := Validate("not a struct")

	if errors == nil {
		t.Fatal("expected errors for non-struct input, got nil")
	}

	if _, ok := errors["_error"]; !ok {
		t.Error("expected '_error' key for non-struct input")
	}
}
