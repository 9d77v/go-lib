package v6

import (
	"context"
	"reflect"
	"testing"

	"github.com/9d77v/go-lib/clients/config"
)

const mapping = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"doc":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
            "retweets":{
                "type":"long"
            },
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

func TestNewClient(t *testing.T) {
	type args struct {
		config *config.ElasticConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{"should be ok", args{&config.ElasticConfig{URLs: []string{"http://localhost:8505/"}}}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CreateIndex(t *testing.T) {
	c, err := NewClient(&config.ElasticConfig{URLs: []string{"http://localhost:8505/"}})
	if err != nil {
		t.Error("连接es失败")
	}
	indexName := c.GetNewIndexName("test", "2006.01.02-15:04")
	type args struct {
		ctx       context.Context
		indexName string
		mapping   string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{"shold be ok", c, args{context.Background(), indexName, mapping}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.CreateIndex(tt.args.ctx, tt.args.indexName, tt.args.mapping); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SetNewAlias(t *testing.T) {
	c, err := NewClient(&config.ElasticConfig{URLs: []string{"http://localhost:8505/"}})
	if err != nil {
		t.Error("连 接es失败 ")
	}
	aliasName := "test"
	layout := "2006.01.02-15:04:05"
	indexName := c.GetNewIndexName(aliasName, layout)
	ctx := context.Background()
	c.CreateIndex(ctx, indexName, mapping)
	type args struct {
		ctx          context.Context
		aliasName    string
		newIndexName string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{"should be ok", c, args{ctx, aliasName, indexName}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetNewAlias(tt.args.ctx, tt.args.aliasName, tt.args.newIndexName); (err != nil) != tt.wantErr {
				t.Errorf("Client.SetNewAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_KeepIndex(t *testing.T) {
	c, err := NewClient(&config.ElasticConfig{URLs: []string{"http://localhost:8505/"}})
	if err != nil {
		t.Error("连接es失败 ")
	}
	ctx := context.Background()
	indexNames := c.FindIndexesByAlias(ctx, "test", "2006.01.02-15:04:05")
	type args struct {
		ctx        context.Context
		indexNames []string
		max        int
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{"should be ok", c, args{ctx, indexNames, 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.KeepIndex(tt.args.ctx, tt.args.indexNames, tt.args.max); (err != nil) != tt.wantErr {
				t.Errorf("Client.SetNewAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
