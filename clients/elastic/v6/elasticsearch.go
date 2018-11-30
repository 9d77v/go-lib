package v6

import (
	"fmt"
	"github.com/9d77v/go-lib/clients/config"
	"github.com/olivere/elastic"
)

//Client .
type Client struct {
	*elastic.Client
}

//NewClient 初始化客户端
func NewClient(config *config.ElasticConfig) (*Client, error) {
	fmt.Println(config.URLs)
	c, err := elastic.NewClient(elastic.SetURL(config.URLs...))
	return &Client{
		c,
	}, err
}

//AggsParam 聚合搜索条件
type AggsParam struct {
	Field string
	Size  int
}

//Aggs 聚合搜索
func Aggs(s *elastic.SearchService, params ...*AggsParam) *elastic.SearchService {
	if s != nil {
		for _, param := range params {
			s = s.Aggregation("group_by_"+param.Field, elastic.NewTermsAggregation().
				Field(param.Field).Size(param.Size))
		}
	}
	return s
}
