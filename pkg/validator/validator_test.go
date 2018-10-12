package validator

import (
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	tcs := []struct {
		form interface{}
		want *Errors
		str  string
	}{
		{
			form: &struct {
				Name string `validate:"required"`
			}{
				Name: "not empty",
			},
			want: nil,
			str:  "",
		},
		{
			form: &struct {
				Name string `validate:"required"`
			}{},
			want: &Errors{
				Data: map[string][]string{
					"name": {"required"},
				},
			},
			str: "name: required;",
		},
	}
	for _, tc := range tcs {
		errs := Validate(tc.form)
		if !reflect.DeepEqual(errs, tc.want) {
			t.Errorf("got %v, want %v\n", errs, tc.want)
		}

		if errs != nil && !reflect.DeepEqual(errs.String(), tc.want.String()) {
			t.Errorf("got %v, want %v\n", errs.String(), tc.want.String())
		}
	}
}

func TestAddError(t *testing.T) {
	tcs := []struct {
		origin     *Errors
		field, msg string
		want       *Errors
		str        string
	}{
		{
			origin: nil,
			field:  "a",
			msg:    "b",
			want:   nil,
			str:    "",
		},
		{
			origin: &Errors{},
			field:  "a",
			msg:    "b",
			want: &Errors{
				Data: map[string][]string{
					"a": {"b"},
				},
			},
			str: "a: b;",
		},
		{
			origin: &Errors{
				Data: map[string][]string{
					"a": {"a"},
				},
			},
			field: "a",
			msg:   "b",
			want: &Errors{
				Data: map[string][]string{
					"a": {"a", "b"},
				},
			},
			str: "a: a,b;",
		},
	}
	for _, tc := range tcs {
		AddError(tc.origin, tc.field, tc.msg)
		if !reflect.DeepEqual(tc.origin, tc.want) {
			t.Errorf("got %v, want %v\n", tc.origin, tc.want)
		}
	}
}
