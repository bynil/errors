package errors

import (
	"errors"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		err  string
		want error
	}{
		{"", fmt.Errorf("")},
		{"foo", fmt.Errorf("foo")},
		{"foo", New("foo")},
		{"string with format specifiers: %v", errors.New("string with format specifiers: %v")},
	}

	for _, tt := range tests {
		got := New(tt.err)
		if got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestWrapNil(t *testing.T) {
	got := Wrap(nil, "no error")
	if got != nil {
		t.Errorf("Wrap(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrap(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := Wrap(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrap(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	x := New("error")
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		// caused error returns cause
		err:  Wrap(io.EOF, "ignored"),
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}, {
		WithMessage(nil, "whoops"),
		nil,
	}, {
		WithMessage(io.EOF, "whoops"),
		io.EOF,
	}, {
		WithStack(nil),
		nil,
	}, {
		WithStack(io.EOF),
		io.EOF,
	}}

	for i, tt := range tests {
		got := Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

func TestWrapfNil(t *testing.T) {
	got := Wrapf(nil, "no error")
	if got != nil {
		t.Errorf("Wrapf(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrapf(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrapf(io.EOF, "read error without format specifiers"), "client error", "client error: read error without format specifiers: EOF"},
		{Wrapf(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		got := Wrapf(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrapf(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errorf("read error without format specifiers"), "read error without format specifiers"},
		{Errorf("read error with %d format specifier", 1), "read error with 1 format specifier"},
	}

	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("Errorf(%v): got: %q, want %q", tt.err, got, tt.want)
		}
	}
}

func TestWithStackNil(t *testing.T) {
	got := WithStack(nil)
	if got != nil {
		t.Errorf("WithStack(nil): got %#v, expected nil", got)
	}
}

func TestWithStack(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{io.EOF, "EOF"},
		{WithStack(io.EOF), "EOF"},
	}

	for _, tt := range tests {
		got := WithStack(tt.err).Error()
		if got != tt.want {
			t.Errorf("WithStack(%v): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestWithMessageNil(t *testing.T) {
	got := WithMessage(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessage(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessage(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := WithMessage(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestWithMessagefNil(t *testing.T) {
	got := WithMessagef(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessagef(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessagef(io.EOF, "read error without format specifier"), "client error", "client error: read error without format specifier: EOF"},
		{WithMessagef(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		got := WithMessagef(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

// errors.New, etc values are not expected to be compared by value
// but the change in errors#27 made them incomparable. Assert that
// various kinds of errors have a functional equality operator, even
// if the result of that equality is always false.
func TestErrorEquality(t *testing.T) {
	vals := []error{
		nil,
		io.EOF,
		errors.New("EOF"),
		New("EOF"),
		Errorf("EOF"),
		Wrap(io.EOF, "EOF"),
		Wrapf(io.EOF, "EOF%d", 2),
		WithMessage(nil, "whoops"),
		WithMessage(io.EOF, "whoops"),
		WithStack(io.EOF),
		WithStack(nil),
	}

	for i := range vals {
		for j := range vals {
			_ = vals[i] == vals[j] // mustn't panic
		}
	}
}

func TestWrapType(t *testing.T) {
	type args struct {
		err     error
		eType   Typer
		message string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantMsg  string
	}{
		{
			name: "wrap type",
			args: args{
				err:     errors.New("test error"),
				eType:   TypeNotFound,
				message: "resource not found",
			},
			wantCode: http.StatusNotFound,
			wantMsg:  "resource not found: test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapType(tt.args.err, tt.args.eType, tt.args.message)
			code, msg := GetAPIError(got)
			if code != tt.wantCode {
				t.Errorf("WrapType() code = %v, want %v", code, tt.wantCode)
			}
			if msg != tt.wantMsg {
				t.Errorf("WrapType() msg = %v, want %v", msg, tt.wantMsg)
			}
		})
	}
}

func TestWrapTypef(t *testing.T) {
	type args struct {
		err    error
		eType  Typer
		format string
		args   []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantMsg  string
	}{
		{
			name: "wrap type",
			args: args{
				err:    errors.New("test error"),
				eType:  TypeNotFound,
				format: "error %s",
				args:   []interface{}{"param"},
			},
			wantCode: http.StatusNotFound,
			wantMsg:  "error param: test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapTypef(tt.args.err, tt.args.eType, tt.args.format, tt.args.args...)
			code, msg := GetAPIError(got)
			if code != tt.wantCode {
				t.Errorf("WrapTypef() code = %v, want %v", code, tt.wantCode)
			}
			if msg != tt.wantMsg {
				t.Errorf("WrapTypef() msg = %v, want %v", msg, tt.wantMsg)
			}
		})
	}
}

func TestLocalization(t *testing.T) {
	oneLc := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "PersonCats",
			One:   "{{.Name}} has {{.Count}} cat.",
			Other: "{{.Name}} has {{.Count}} cats.",
		},
		TemplateData: map[string]interface{}{
			"Name":  "Nick",
			"Count": 1,
		},
		PluralCount: 1,
	}
	twoLc := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "PersonCats",
			One:   "{{.Name}} has {{.Count}} cat.",
			Other: "{{.Name}} has {{.Count}} cats.",
		},
		TemplateData: map[string]interface{}{
			"Name":  "Nick",
			"Count": 2,
		},
		PluralCount: 2,
	}
	tests := []struct {
		eType    Typer
		lc       *i18n.LocalizeConfig
		want     string
		wantType Typer
	}{
		{TypeDuplicate, oneLc, "Nick has 1 cat.", TypeDuplicate},
		{TypeNotFound, twoLc, "Nick has 2 cats.", TypeNotFound},
	}

	for _, tt := range tests {
		got := NewI18n(tt.eType, tt.lc).Error()
		if got != tt.want {
			t.Errorf("NewI18n(%v, %v): got: %q, want %q", tt.eType, tt.lc, got, tt.want)
		}
		gotLc := NewI18n(tt.eType, tt.lc).(I18ner).LocalizeConfig()
		if gotLc != tt.lc {
			t.Errorf("NewI18n(%v, %v): got: %v, want %v", tt.eType, tt.lc, gotLc, tt.lc)
		}
		if !HasType(NewI18n(tt.eType, tt.lc), tt.wantType) {
			t.Errorf("NewI18n(%v, %v): do not have type %v", tt.eType, tt.lc, tt.wantType)
		}
	}
}
