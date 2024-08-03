package logger

import "context"

type key string

const (
	userUIDKey   key = "logger.user_uid"
	requestIDKey key = "logger.request_id"
)

// AddUserUID return context with tenantUID
func AddUserUID(ctx context.Context, userUID string) context.Context {
	return context.WithValue(ctx, userUIDKey, userUID)
}

// AddRequestID return context with RequestID
func AddRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, reqID)
}

func extractParams(ctx context.Context) (args []any) {
	if userUID, ok := ctx.Value(userUIDKey).(string); ok && userUID != "" {
		args = append(args, "tenant_uid", userUID)
	}
	if reqID, ok := ctx.Value(requestIDKey).(string); ok && reqID != "" {
		args = append(args, "request_id", reqID)
	}

	return args
}
