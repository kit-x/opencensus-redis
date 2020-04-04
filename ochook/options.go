package ochook

import (
	"github.com/go-redis/redis/v7"
	"go.opencensus.io/trace"
)

type TraceOption func(o *TraceOptions)

type TraceOptions struct {
	// AllowRoot, if set to true, will allow hook to create root spans in
	// absence of existing spans or even context.
	// Default is to not trace hook calls if no existing parent span is found
	// in context or when using methods not taking context.
	AllowRoot bool

	// Enable decide which cmd should be traced
	// Default allow all cmd expect "ping"
	Enable Decider

	// DefaultAttributes will be set to each span as default.
	DefaultAttributes []trace.Attribute
}

var _defaultOptions = TraceOptions{
	Enable: func(cmd redis.Cmder) bool {
		return cmd.Name() != "ping"
	},
}

type Decider func(cmd redis.Cmder) bool

// WithPing if set to true, will enable the creation of spans on Ping requests.
func WithDecider(fn Decider) TraceOption {
	return func(o *TraceOptions) {
		o.Enable = fn
	}
}

// WithDefaultAttributes will be set to each span as default.
func WithDefaultAttributes(attrs ...trace.Attribute) TraceOption {
	return func(o *TraceOptions) {
		o.DefaultAttributes = attrs
	}
}

// WithAllowRoot if set to true, will allow ocsql to create root spans in
// absence of exisiting spans or even context.
// Default is to not trace ocsql calls if no existing parent span is found
// in context or when using methods not taking context.
func WithAllowRoot(b bool) TraceOption {
	return func(o *TraceOptions) {
		o.AllowRoot = b
	}
}
