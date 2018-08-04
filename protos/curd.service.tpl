package services

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"

	pb "github.com/9d77v/${project}-${module}/protos"
)

//${Entity}Service ...
type ${Entity}Service struct{}

//Create${Entity} ...
func (s *${Entity}Service) Create${Entity}(ctx context.Context,
	in *pb.Create${Entity}Request) (response *pb.Create${Entity}Response, err error) {
	response = new(pb.Create${Entity}Response)
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}Service.Create${Entity}", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCServer.Set(span)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	return response, nil
}

//Update${Entity} ...
func (s *${Entity}Service) Update${Entity}(ctx context.Context,
	in *pb.Update${Entity}Request) (response *pb.Update${Entity}Response, err error) {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}Service.Update${Entity}", opentracing.ChildOf(span.Context()))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return response, nil
}

//Delete${Entity} ...
func (s *${Entity}Service) Delete${Entity}(ctx context.Context,
	in *pb.Delete${Entity}Request) (response *pb.Delete${Entity}Response, err error) {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}Service.Delete${Entity}", opentracing.ChildOf(span.Context()))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	return response, nil
}

//Get${Entity}ByID ...
func (s *${Entity}Service) Get${Entity}ByID(ctx context.Context,
	in *pb.Get${Entity}ByIDRequest) (response *pb.Get${Entity}ByIDResponse, err error) {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}Service.Get${Entity}ByID", opentracing.ChildOf(span.Context()))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	return response, nil
}

//List${Entity} ...
func (s *${Entity}Service) List${Entity}(ctx context.Context,
	in *pb.List${Entity}Request) (response *pb.List${Entity}Response, err error) {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("${Entity}Service.List${Entity}", opentracing.ChildOf(span.Context()))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	return response, nil
}