package ochook

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v7"
	"go.opencensus.io/trace"
)

type redisHook struct {
	opt *TraceOptions
}

func New(opts ...TraceOption) redis.Hook {
	opt := _defaultOptions
	for _, o := range opts {
		o(&opt)
	}
	return &redisHook{
		opt: &opt,
	}
}

func (r *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !r.opt.Enable(cmd) {
		return ctx, nil
	}
	if r.opt.AllowRoot || trace.FromContext(ctx) != nil {
		var span *trace.Span
		ctx, span = trace.StartSpan(ctx, "redis.process."+cmd.Name(), trace.WithSpanKind(trace.SpanKindClient))
		if len(r.opt.DefaultAttributes) > 0 {
			span.AddAttributes(r.opt.DefaultAttributes...)
		}
		span.AddAttributes(
			trace.StringAttribute("cmd", cmd.Name()),
			trace.StringAttribute("args", fmt.Sprintf("%v", cmd.Args())),
		)
	}

	return ctx, nil
}

func (r *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if !r.opt.Enable(cmd) {
		return nil
	}
	if span := trace.FromContext(ctx); span != nil {
		setSpanStatus(span, cmd.Err())
		span.End()
	}
	return nil
}

func (r *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if r.opt.AllowRoot || trace.FromContext(ctx) != nil {
		var span *trace.Span
		ctx, span = trace.StartSpan(ctx, "redis.pipeline", trace.WithSpanKind(trace.SpanKindClient))
		span.Annotate(attributesFromCommands(cmds), "log")
	}
	return ctx, nil
}

func (r *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if span := trace.FromContext(ctx); span != nil {
		setSpanStatus(span, firstCmdsErr(cmds))
		span.End()
	}
	return nil
}
