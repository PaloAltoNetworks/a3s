package crud

import (
	"context"
	"reflect"
	"testing"

	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

func TestTranslateContext(t *testing.T) {
	type args struct {
		bctx bahamut.Context
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      manipulate.Context
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			"empty request",
			func(t *testing.T) args {
				bctx := bahamut.NewContext(context.Background(), elemental.NewRequest())
				return args{bctx}
			},
			manipulate.NewContext(context.Background()),
			false,
			nil,
		},
		{
			"request with namespace",
			func(t *testing.T) args {
				req := &elemental.Request{
					Namespace: "/test",
				}
				bctx := bahamut.NewContext(context.Background(), req)
				return args{bctx}
			},
			manipulate.NewContext(
				context.Background(),
				manipulate.ContextOptionNamespace("/test"),
			),
			false,
			nil,
		},
		{
			"request with namespace and recursive",
			func(t *testing.T) args {
				req := &elemental.Request{
					Recursive: true,
					Namespace: "/test",
				}
				bctx := bahamut.NewContext(context.Background(), req)
				return args{bctx}
			},
			manipulate.NewContext(
				context.Background(),
				manipulate.ContextOptionNamespace("/test"),
				manipulate.ContextOptionRecursive(true),
			),
			false,
			nil,
		},
		{
			"request with namespace propapated and recursive",
			func(t *testing.T) args {
				req := &elemental.Request{
					Recursive:  true,
					Namespace:  "/test",
					Propagated: true,
				}
				bctx := bahamut.NewContext(context.Background(), req)
				return args{bctx}
			},
			manipulate.NewContext(
				context.Background(),
				manipulate.ContextOptionNamespace("/test"),
				manipulate.ContextOptionRecursive(true),
				manipulate.ContextOptionPropagated(true),
			),
			false,
			nil,
		},
		{
			"request with namespace propapated and recursive and valid filter",
			func(t *testing.T) args {
				req := &elemental.Request{
					Recursive:  true,
					Namespace:  "/test",
					Propagated: true,
					Parameters: elemental.Parameters{
						"q": elemental.NewParameter(elemental.ParameterTypeString, `name == "blah"`),
					},
				}
				bctx := bahamut.NewContext(context.Background(), req)
				return args{bctx}
			},
			manipulate.NewContext(
				context.Background(),
				manipulate.ContextOptionNamespace("/test"),
				manipulate.ContextOptionRecursive(true),
				manipulate.ContextOptionPropagated(true),
				manipulate.ContextOptionFilter(
					elemental.NewFilterComposer().WithKey("name").Equals("blah").Done(),
				),
			),
			false,
			nil,
		},
		{
			"request with invalid filter",
			func(t *testing.T) args {
				req := &elemental.Request{
					Parameters: elemental.Parameters{
						"q": elemental.NewParameter(elemental.ParameterTypeString, `oh noes`),
					},
				}
				bctx := bahamut.NewContext(context.Background(), req)
				return args{bctx}
			},
			nil,
			true,
			func(err error, t *testing.T) {
				if err.Error() != "error 400 (a3s:policy): Bad Request: unable to parse filter in query parameter: invalid operator. found noes instead of (==, !=, <, <=, >, >=, contains, in, matches, exists)" {
					t.Fail()
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := TranslateContext(tArgs.bctx)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("TranslateContext got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("TranslateContext error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
