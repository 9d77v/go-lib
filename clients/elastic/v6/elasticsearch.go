package v6

import (
	"context"
	"fmt"
	"github.com/9d77v/go-lib/clients/config"
	"github.com/9d77v/go-lib/utils"
	"github.com/9d77v/go-lib/worker"
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

//BulkDoc es doc
type BulkDoc struct {
	ID  string
	Doc interface{}
}

//BulkInsert 多worker按一定数量批量导入es
func (esClient *Client) BulkInsert(ctx context.Context, bds []*BulkDoc, indexName string, bulkNum, workerNum int) []error {
	size := len(bds)
	max := size/bulkNum + 1
	pool := worker.NewStaticPool(workerNum)
	errs := make([]error, 0, max)
	for i := 0; i < max; i++ {
		i := i
		pool.Add(func() {
			bulkService := elastic.NewBulkService(esClient.Client)
			for j := i * bulkNum; j < utils.MinInt(size, (j+1)*bulkNum); j++ {
				req := elastic.NewBulkIndexRequest()
				req.Index(indexName).
					Type("doc").
					Id(bds[j].ID).
					Doc(bds[j].Doc)
				bulkService.Add(req)
			}
			_, err := bulkService.Do(ctx)
			if err != nil {
				errs = append(errs, err)
			}
		})
	}
	pool.Stop()
	return errs
}
