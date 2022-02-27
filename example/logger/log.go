package main

import (
	"context"

	"github.com/aivencs/kit/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	logger.Init("Example", "product", "json")
	ctx := context.WithValue(context.Background(), "trace", "109873")
	logger.Info(ctx, "example", zap.Any("param", "aivenc"))
	otherExample()
}

func otherExample() {
	ctx := context.WithValue(context.Background(), "trace", "87000")
	logger.Error(ctx, "aivenc", zap.Any("param", "example"))
}
