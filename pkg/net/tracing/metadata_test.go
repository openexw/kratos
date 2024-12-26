package tracing

import (
	"context"
	"kratos/pkg/net/metadata"
	"os"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/propagation"
)

func TestMain(m *testing.M) {
	os.Setenv("APP_NAME", "testDemo")
	m.Run()
}

func TestMetadata_Inject(t *testing.T) {
	type args struct {
		appName string
		carrier propagation.TextMapCarrier
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "https://go-kratos.dev",
			args: args{"https://go-kratos.dev", propagation.HeaderCarrier{}},
			want: "https://go-kratos.dev",
		},
		{
			name: "https://github.com/go-kratos/kratos",
			args: args{"https://github.com/go-kratos/kratos", propagation.HeaderCarrier{"mode": []string{"test"}}},
			want: "https://github.com/go-kratos/kratos",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("APP_NAME", tt.name)
			ctx := context.Background()
			m := new(Metadata)
			m.Inject(ctx, tt.args.carrier)
			if res := tt.args.carrier.Get(serviceHeader); tt.want != res {
				t.Errorf("Get(serviceHeader) :%s want: %s", res, tt.want)
			}
		})
	}
}

func TestMetadata_Extract(t *testing.T) {
	type args struct {
		parent  context.Context
		carrier propagation.TextMapCarrier
	}
	tests := []struct {
		name  string
		args  args
		want  string
		crash bool
	}{
		{
			name: "https://go-kratos.dev",
			args: args{
				parent:  context.Background(),
				carrier: propagation.HeaderCarrier{"X-Md-Service-Name": []string{"https://go-kratos.dev"}},
			},
			want: "https://go-kratos.dev",
		},
		{
			name: "https://github.com/go-kratos/kratos",
			args: args{
				parent:  metadata.NewServerContext(context.Background(), metadata.MD{}),
				carrier: propagation.HeaderCarrier{"X-Md-Service-Name": []string{"https://github.com/go-kratos/kratos"}},
			},
			want: "https://github.com/go-kratos/kratos",
		},
		{
			name: "https://github.com/go-kratos/kratos",
			args: args{
				parent:  metadata.NewServerContext(context.Background(), metadata.MD{}),
				carrier: propagation.HeaderCarrier{"X-Md-Service-Name": nil},
			},
			crash: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Metadata{}
			os.Setenv("APP_NAME", tt.name)
			ctx := b.Extract(tt.args.parent, tt.args.carrier)
			md, ok := metadata.FromServerContext(ctx)
			if !ok {
				if tt.crash {
					return
				}
				t.Errorf("expect %v, got %v", true, ok)
			}
			if !reflect.DeepEqual(md.Get(serviceHeader), tt.want) {
				t.Errorf("expect %v, got %v", tt.want, md.Get(serviceHeader))
			}
		})
	}
}

func TestFields(t *testing.T) {
	b := Metadata{}
	if !reflect.DeepEqual(b.Fields(), []string{"x-md-service-name"}) {
		t.Errorf("expect %v, got %v", []string{"x-md-service-name"}, b.Fields())
	}
}
