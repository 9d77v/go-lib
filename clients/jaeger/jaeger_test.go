package jaeger

import (
	"testing"
	"time"

	"github.com/9d77v/go-lib/clients/etcd"
)

func TestInitTracerFromEtcd(t *testing.T) {
	etcdCli, err := etcd.NewClient(5 * time.Second)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	type args struct {
		etcdCli     *etcd.Client
		serviceName string
	}
	tests := []struct {
		name string
		args args
		// wantTracer opentracing.Tracer
		// wantCloser io.Closer
		wantErr bool
	}{
		{"should be ok", args{etcdCli, "redis"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := InitTracerFromEtcd(tt.args.etcdCli, tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitTracerFromEtcd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotTracer, tt.wantTracer) {
			// 	t.Errorf("InitTracerFromEtcd() gotTracer = %v, want %v", gotTracer, tt.wantTracer)
			// }
			// if !reflect.DeepEqual(gotCloser, tt.wantCloser) {
			// 	t.Errorf("InitTracerFromEtcd() gotCloser = %v, want %v", gotCloser, tt.wantCloser)
			// }
		})
	}
	time.Sleep(1 * time.Second)
}

func TestInitGlobalTracerFromEtcd(t *testing.T) {
	etcdCli, err := etcd.NewClient(5 * time.Second)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	type args struct {
		etcdCli *etcd.Client
	}
	tests := []struct {
		name string
		args args
		// wantCloser io.Closer
		wantErr bool
	}{
		{"should be ok", args{etcdCli}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := InitGlobalTracerFromEtcd(tt.args.etcdCli)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitGlobalTracerFromEtcd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotCloser, tt.wantCloser) {
			// 	t.Errorf("InitGlobalTracerFromEtcd() = %v, want %v", gotCloser, tt.wantCloser)
			// }
		})
	}
	time.Sleep(1 * time.Second)
}
