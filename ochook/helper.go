package ochook

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v7"
	"go.opencensus.io/trace"
)

func attributesFromCommands(cmds []redis.Cmder) []trace.Attribute {
	if len(cmds) == 0 {
		return nil
	}
	attributes := make([]trace.Attribute, len(cmds))
	for i, cmd := range cmds {
		attributes[i] = trace.StringAttribute(fmt.Sprintf("cmd_%d", i), fmt.Sprintf("%v", cmd.Args()))
	}
	return attributes
}

func firstCmdsErr(cmds []redis.Cmder) error {
	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			return err
		}
	}
	return nil
}

func setSpanStatus(span *trace.Span, err error) {
	var status trace.Status
	switch err {
	case nil:
		status.Code = trace.StatusCodeOK
		span.SetStatus(status)
		return
	case redis.Nil:
		status.Code = trace.StatusCodeNotFound
	case redis.TxFailedErr:
		status.Code = trace.StatusCodeFailedPrecondition
	case context.Canceled:
		status.Code = trace.StatusCodeCancelled
	case context.DeadlineExceeded:
		status.Code = trace.StatusCodeDeadlineExceeded
	default:
		status.Code = trace.StatusCodeUnknown
	}
	status.Message = err.Error()
	span.SetStatus(status)
}
