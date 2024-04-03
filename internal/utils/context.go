package utils

import "context"

const (
	FlowCtx = "FLOW_CTX"
)

func GetStringFromContext(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}
	stringValue, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return stringValue
}

func GetContextWithFlowName(ctx context.Context, flowName string) context.Context {
	return context.WithValue(ctx, FlowCtx, flowName)
}
