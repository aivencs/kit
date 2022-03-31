package main

import (
	"context"
	"fmt"

	"github.com/aivencs/kit/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	logger.InitErc()
	errorc := logger.GetDefaultErc()
	fmt.Println(errorc)
	logger.InitLogger("zap", "service-work", "product", "label-name", "json")
	ctx := context.WithValue(context.Background(), "trace", "109873")
	logger.Info(ctx, "example", zap.Any("param", "aivenc"))
	otherExample()
}

func otherExample() {
	ctx := context.WithValue(context.Background(), "trace", "87000")
	logger.Error(ctx, "aivenc", zap.Any("param", "example"))
}
