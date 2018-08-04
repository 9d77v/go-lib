package gorm

import (
	"testing"
	"time"

	"github.com/9d77v/go-lib/clients/etcd"

	"github.com/9d77v/go-lib/clients/config"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestNewGormClient(t *testing.T) {
	cli, err := etcd.NewClient(5 * time.Second)
	dbConfig := new(config.DBConfig)
	dbKey := "/config/dev/eshop-order/db"
	err = cli.GetValue(10*time.Second, dbKey, dbConfig)
	if err != nil {
		t.Error("db config not exist:", err)
	}
	type args struct {
		config *config.DBConfig
	}
	tests := []struct {
		name string
		args args
		// want    *GormClient
		wantErr bool
	}{
		{"TestNewGormClient", args{dbConfig}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGormClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewGormClient() = %v, want %v", got, tt.want)
			// }
		})
	}
}
