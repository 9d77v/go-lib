package gorm

import (
	"testing"
	"time"

	"github.com/9d77v/go-lib/clients/etcd"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestNewClientFromEtcd(t *testing.T) {
	etcdCli, err := etcd.NewClient(5 * time.Second)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	type args struct {
		etcdCli *etcd.Client
		values  []interface{}
	}
	tests := []struct {
		name string
		args args
		// wantDbCli *Client
		wantErr bool
	}{
		{"should be ok", args{etcdCli, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClientFromEtcd(tt.args.etcdCli, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientFromEtcd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotDbCli, tt.wantDbCli) {
			// 	t.Errorf("NewClientFromEtcd() = %v, want %v", gotDbCli, tt.wantDbCli)
			// }
		})
	}
	time.Sleep(25 * time.Second)
}
