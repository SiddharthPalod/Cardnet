package network

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	maxRetries = 1
	retryDelay = 50 * time.Millisecond
)

func Retry(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err = fn(ctx)
		if err == nil {
			return nil
		}

		if !isRetryable(err) {
			return err
		}

		time.Sleep(retryDelay)
	}
	return err
}

func isRetryable(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return true
	}

	switch st.Code() {
	case codes.DeadlineExceeded,
		codes.Unavailable:
		return true
	default:
		return false
	}
}
