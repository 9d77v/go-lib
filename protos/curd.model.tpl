package models

import (
	"context"
	"time"

	"github.com/lib/pq"
	opentracing "github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
)

//${Entity} ...
type ${Entity} struct {
	ID          int64
	CreateTime  time.Time     `gorm:"NOT NULL"`
	UpdateTime  time.Time     `gorm:"NOT NULL"`
}

//New${Entity} ...
func New${Entity}( id ...int64) *${Entity} {
	c := &${Entity}{
	}
	if id != nil {
		c.ID = id[0]
	}
	return c
}

//Create create a new model
func (c *${Entity}) Create(ctx context.Context) error {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}.Create", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCServer.Set(span)
		tags.PeerService.Set(span, "postgres")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	db := Gorm.Create(c)
	if span != nil {
		span.LogKV("${Entity}", c)
		span.LogKV("db.Error", db.Error)
	}
	return db.Error
}

//Update update model
func (c *${Entity}) Update(ctx context.Context, new${Entity} *${Entity}) error {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}.Update", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCServer.Set(span)
		tags.PeerService.Set(span, "postgres")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	tx := Gorm.Begin()
	if span != nil {
		defer span.LogKV("${Entity}", c)
		defer span.LogKV("db.Error", tx.Error)
	}
	if err := tx.First(c, new${Entity}.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(c).Updates(map[string]interface{}{
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//Delete delete model
func (c *${Entity}) Delete(ctx context.Context, id int64) error {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}.Delete", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCServer.Set(span)
		tags.PeerService.Set(span, "postgres")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	tx := Gorm.Begin()
	if span != nil {
		defer span.LogKV("${Entity}", c)
		defer span.LogKV("db.Error", tx.Error)
	}
	if err := tx.First(c, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(c).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//GetByID delete model
func (c *${Entity}) GetByID(ctx context.Context, id int64) error {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}.GetByID", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCServer.Set(span)
		tags.PeerService.Set(span, "postgres")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	db := Gorm.First(c, id)
	if span != nil {
		span.LogKV("id", id)
		span.LogKV("${Entity}", c)
		span.LogKV("db.Error", db.Error)
	}
	return db.Error
}

//List delete model
func (c *${Entity}) List(ctx context.Context,
	offset, limit int64) (list []*${Entity}, total int64, err error) {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}.List", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCServer.Set(span)
		tags.PeerService.Set(span, "postgres")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	if span != nil {
		span.LogKV("offset", offset)
		span.LogKV("limit", limit)
		defer span.LogKV("db.Error", err)
	}
	list = make([]*${Entity}, 0)
	if err = Gorm.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	if err = Gorm.Model(c).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
