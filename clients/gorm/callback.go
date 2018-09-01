package gorm

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	opentracing "github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
)

type callbacks struct{}

func newCallbacks() *callbacks {
	return &callbacks{}
}

func (c *callbacks) beforeCreate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterCreate(scope *gorm.Scope)    { c.after(scope, "INSERT") }
func (c *callbacks) beforeQuery(scope *gorm.Scope)    { c.before(scope) }
func (c *callbacks) afterQuery(scope *gorm.Scope)     { c.after(scope, "SELECT") }
func (c *callbacks) beforeUpdate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterUpdate(scope *gorm.Scope)    { c.after(scope, "UPDATE") }
func (c *callbacks) beforeDelete(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterDelete(scope *gorm.Scope)    { c.after(scope, "DELETE") }
func (c *callbacks) beforeRowQuery(scope *gorm.Scope) { c.before(scope) }
func (c *callbacks) afterRowQuery(scope *gorm.Scope)  { c.after(scope, "") }

func (c *callbacks) before(scope *gorm.Scope) {
	parentSpanValue, ok := scope.Get(parentSpanGormKey)
	if !ok {
		return
	}
	tracerValue, ok := scope.Get(spanGormTracer)
	if !ok {
		return
	}
	parentSpan := parentSpanValue.(opentracing.Span)
	tr := tracerValue.(opentracing.Tracer)
	span := tr.StartSpan("sql", opentracing.ChildOf(parentSpan.Context()))
	tags.DBType.Set(span, "sql")
	scope.Set(spanGormKey, span)
}

func (c *callbacks) after(scope *gorm.Scope, operation string) {
	if scope.HasError() {
		return
	}
	val, ok := scope.Get(spanGormKey)
	if !ok {
		return
	}
	span := val.(opentracing.Span)
	if operation == "" {
		operation = strings.ToUpper(strings.Split(scope.SQL, " ")[0])
	}
	tags.DBStatement.Set(span, fmt.Sprint(gorm.LogFormatter(scope.SQL, scope.SQLVars)...))
	span.SetTag("db.table", scope.TableName())
	span.SetTag("db.method", operation)
	span.Finish()
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)
	// gorm does some magic, if you pass CallbackProcessor here - nothing works
	switch name {
	case "create":
		db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate)
		db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)
	case "query":
		db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery)
		db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)
	case "update":
		db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate)
		db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)
	case "delete":
		db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete)
		db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)
	case "row_query":
		db.Callback().RowQuery().Before(gormCallbackName).Register(beforeName, c.beforeRowQuery)
		db.Callback().RowQuery().After(gormCallbackName).Register(afterName, c.afterRowQuery)
	}
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		now := time.Now()

		if createdAtField, ok := scope.FieldByName("CreateTime"); ok {
			if createdAtField.IsBlank {
				createdAtField.Set(now)
			}
		}

		if updatedAtField, ok := scope.FieldByName("UpdateTime"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(now)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdateTime", time.Now())
	}
}
