package temporal

import (
	"context"

	"go.fabra.io/server/common/errors"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/temporal"
)

type workerInterceptor struct {
	interceptor.WorkerInterceptorBase
}

func NewErrorInterceptor() interceptor.WorkerInterceptor {
	return &workerInterceptor{}
}

func (w *workerInterceptor) InterceptActivity(
	ctx context.Context,
	next interceptor.ActivityInboundInterceptor,
) interceptor.ActivityInboundInterceptor {
	i := &activityInboundInterceptor{}
	i.Next = next
	return i
}

type activityInboundInterceptor struct {
	interceptor.ActivityInboundInterceptorBase
}

func (a *activityInboundInterceptor) ExecuteActivity(
	ctx context.Context,
	in *interceptor.ExecuteActivityInput,
) (interface{}, error) {
	result, err := a.Next.ExecuteActivity(ctx, in)

	// Check if the error is a CustomerVisibleError and return that as the top-level error
	// We don't want to expose any other information that might have been added due to wrapping
	var customerVisisbleError *errors.CustomerVisibleError
	if errors.As(err, &customerVisisbleError) {
		return result, temporal.NewApplicationErrorWithCause(customerVisisbleError.Error(), "CustomerVisibleError", err)
	} else {
		return result, err
	}
}
