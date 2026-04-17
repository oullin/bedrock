package crm

import (
	"strings"
	"testing"
)

func TestContactFormValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		form contactForm
		want map[string]string
	}{
		{
			name: "missing required fields",
			form: contactForm{},
			want: map[string]string{
				"first_name": "The first name field is required.",
				"last_name":  "The last name field is required.",
				"email":      "The email field is required.",
			},
		},
		{
			name: "invalid email",
			form: contactForm{
				FirstName: "Mina",
				LastName:  "Cole",
				Email:     "mina.example.test",
			},
			want: map[string]string{
				"email": "The email field must be a valid email address.",
			},
		},
		{
			name: "first_name exceeds max length",
			form: contactForm{
				FirstName: strings.Repeat("a", 256),
				LastName:  "Cole",
				Email:     "mina@example.test",
			},
			want: map[string]string{
				"first_name": "The first name field must not exceed 255 characters.",
			},
		},
		{
			name: "phone exceeds max length",
			form: contactForm{
				FirstName: "Mina",
				LastName:  "Cole",
				Email:     "mina@example.test",
				Phone:     strings.Repeat("5", 256),
			},
			want: map[string]string{
				"phone": "The phone field must not exceed 255 characters.",
			},
		},
		{
			name: "non-numeric organization_id",
			form: contactForm{
				FirstName:      "Mina",
				LastName:       "Cole",
				Email:          "mina@example.test",
				OrganizationID: "abc",
			},
			want: map[string]string{
				"organization_id": "The organization id field must be a valid identifier.",
			},
		},
		{
			name: "decimal organization_id",
			form: contactForm{
				FirstName:      "Mina",
				LastName:       "Cole",
				Email:          "mina@example.test",
				OrganizationID: "12.34",
			},
			want: map[string]string{
				"organization_id": "The organization id field must be a valid identifier.",
			},
		},
		{
			name: "valid form with organization_id",
			form: contactForm{
				FirstName:      "Mina",
				LastName:       "Cole",
				Email:          "mina@example.test",
				OrganizationID: "42",
			},
			want: map[string]string{},
		},
		{
			name: "valid form",
			form: contactForm{
				FirstName: "Mina",
				LastName:  "Cole",
				Email:     "mina@example.test",
			},
			want: map[string]string{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.form.validate()

			if len(got) != len(tt.want) {
				t.Fatalf("len(validate()) = %d, want %d; got=%v", len(got), len(tt.want), got)
			}

			for key, want := range tt.want {
				if got[key] != want {
					t.Fatalf("validate()[%q] = %q, want %q", key, got[key], want)
				}
			}
		})
	}
}

func TestOrganizationFormValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		form organizationForm
		want map[string]string
	}{
		{
			name: "missing required fields",
			form: organizationForm{},
			want: map[string]string{
				"name": "The name field is required.",
			},
		},
		{
			name: "name exceeds max length",
			form: organizationForm{Name: strings.Repeat("a", 256)},
			want: map[string]string{
				"name": "The name field must not exceed 255 characters.",
			},
		},
		{
			name: "valid form",
			form: organizationForm{Name: "Acme"},
			want: map[string]string{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.form.validate()

			if len(got) != len(tt.want) {
				t.Fatalf("len(validate()) = %d, want %d; got=%v", len(got), len(tt.want), got)
			}

			for key, want := range tt.want {
				if got[key] != want {
					t.Fatalf("validate()[%q] = %q, want %q", key, got[key], want)
				}
			}
		})
	}
}
