package errors

import (
	"errors"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
	"reflect"
	"testing"
)

type CustomTypeExample int

func (e CustomTypeExample) HTTPStatusCode() int { return int(e) }

func TestHasType(t *testing.T) {
	type args struct {
		err error
		et  Typer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "has required type, nested",
			args: args{
				err: WrapType(
					WrapType(
						Internal("hello world"),
						TypeEmpty,
						DefaultMessage,
					),
					TypeUnauthenticated,
					DefaultMessage,
				),
				et: TypeInternal,
			},
			want: true,
		},
		{
			name: "has required type, not nested",
			args: args{
				err: Internal("hello world"),
				et:  TypeInternal,
			},
			want: true,
		},
		{
			name: "has required custom type, nested",
			args: args{
				err: WrapType(
					WrapType(
						Internal("hello world"),
						CustomTypeExample(499),
						DefaultMessage,
					),
					TypeUnauthenticated,
					DefaultMessage,
				),
				et: CustomTypeExample(499),
			},
			want: true,
		},
		{
			name: "does not have required type",
			args: args{
				err: WrapType(
					WrapType(
						Internal("hello world"),
						TypeEmpty,
						DefaultMessage,
					),
					TypeUnauthenticated,
					DefaultMessage,
				),
				et: TypeValidation,
			},
			want: false,
		},
		{
			name: "*Error wrapped in external error",
			args: args{
				err: fmt.Errorf("unknown error %w", Internal("internal error")),
				et:  TypeInternal,
			},
			want: true,
		},
		{
			name: "other error type",
			args: args{
				err: fmt.Errorf("external error"),
				et:  TypeInput,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasType(tt.args.err, tt.args.et); got != tt.want {
				t.Errorf("HasType() = %v, want %v %s", got, tt.want, tt.args.err.Error())
			}
		})
	}
}

func TestGetAPIError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantMsg  string
	}{
		{
			name: "nested errors",
			args: args{
				err: WrapType(
					WrapType(
						Internal("hello world"),
						TypeEmpty,
						DefaultMessage,
					),
					TypeUnauthenticated,
					DefaultMessage,
				),
			},
			wantCode: TypeUnauthenticated.HTTPStatusCode(),
			wantMsg:  DefaultMessage + ": " + DefaultMessage + ": hello world",
		},
		{
			name: "nested normal errors",
			args: args{
				err: fmt.Errorf("unknown error %w", NotFound("hello world")),
			},
			wantCode: TypeNotFound.HTTPStatusCode(),
			wantMsg:  "unknown error hello world",
		},
		{
			name: "normal errors",
			args: args{
				err: fmt.Errorf("unknown error %w", fmt.Errorf("hello world")),
			},
			wantCode: defaultErrType.HTTPStatusCode(),
			wantMsg:  "unknown error hello world",
		},
		{
			name: "custom type",
			args: args{
				err: WrapType(
					WrapType(
						Internal("hello world"),
						TypeEmpty,
						DefaultMessage,
					),
					CustomTypeExample(499),
					DefaultMessage,
				),
			},
			wantCode: 499,
			wantMsg:  DefaultMessage + ": " + DefaultMessage + ": hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotMsg := GetAPIError(tt.args.err)
			if gotCode != tt.wantCode {
				t.Errorf("GetAPIError() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("GetAPIError() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}

			fmt.Printf("%+v\n", tt.args.err)
		})
	}
}

func TestAllGetAPIError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want2 string
	}{
		{
			name: "TypeInternal",
			args: args{
				err: Internal("unknown error occurred"),
			},
			want:  http.StatusInternalServerError,
			want2: "unknown error occurred",
		},
		{
			name: "TypeInternal - Go builtin error type",
			args: args{
				err: errors.New("unknown error occurred"),
			},
			want:  http.StatusInternalServerError,
			want2: "unknown error occurred",
		},
		{
			name: "TypeValidation",
			args: args{
				err: Validation("invalid email provided"),
			},
			want:  http.StatusUnprocessableEntity,
			want2: "invalid email provided",
		},
		{
			name: "TypeInput",
			args: args{
				err: Input("invalid json provided"),
			},
			want:  http.StatusBadRequest,
			want2: "invalid json provided",
		},
		{
			name: "TypeDuplicate",
			args: args{
				err: Duplicate("duplicate content detected"),
			},
			want:  http.StatusConflict,
			want2: "duplicate content detected",
		},
		{
			name: "TypeUnauthenticated",
			args: args{
				err: Unauthenticated("authentication required"),
			},
			want:  http.StatusUnauthorized,
			want2: "authentication required",
		},
		{
			name: "TypeNoPermission",
			args: args{
				err: NoPermission("not authorized to access this resource"),
			},
			want:  http.StatusForbidden,
			want2: "not authorized to access this resource",
		},
		{
			name: "TypeEmpty",
			args: args{
				err: Empty("empty content not expected"),
			},
			want:  http.StatusGone,
			want2: "empty content not expected",
		},
		{
			name: "TypeNotFound",
			args: args{
				err: NotFound("requested resource not found"),
			},
			want:  http.StatusNotFound,
			want2: "requested resource not found",
		},
		{
			name: "TypeLimitExceeded",
			args: args{
				err: LimitExceeded("exceeded maximum number of requests allowed"),
			},
			want:  http.StatusTooManyRequests,
			want2: "exceeded maximum number of requests allowed",
		},
		{
			name: "TypeSubscriptionExpired",
			args: args{
				err: SubscriptionExpired("your subscription has expired"),
			},
			want:  http.StatusPaymentRequired,
			want2: "your subscription has expired",
		},
		{
			name: "Custom Type",
			args: args{
				err: WrapType(New("internal error"), NewCustomType("error detail", http.StatusFailedDependency), DefaultMessage),
			},
			want:  http.StatusFailedDependency,
			want2: DefaultMessage + ": internal error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := GetAPIError(tt.args.err)
			if got != tt.want {
				t.Errorf("GetAPIError(), %s, got = %v, want %v", tt.name, got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("GetAPIError(), %s, got2 = %v, want %v", tt.name, got2, tt.want2)
			}
		})
	}
}

func TestGetLocalizeConfig(t *testing.T) {
	lc := &i18n.LocalizeConfig{
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
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		args   args
		wantLc *i18n.LocalizeConfig
	}{
		{
			"ok",
			args{err: NewI18n(TypeNotFound, lc)},
			lc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLc := GetLocalizeConfig(tt.args.err); !reflect.DeepEqual(gotLc, tt.wantLc) {
				t.Errorf("GetLocalizeConfig() = %v, want %v", gotLc, tt.wantLc)
			}
		})
	}
}
