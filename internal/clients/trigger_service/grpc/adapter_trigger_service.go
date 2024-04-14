package trigger_hook_grpc

import (
	"context"
	"fmt"
	trg_hk "github.com/fishmanDK/proto_avito_test_task/protos/gen/go/trigger_hook"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	api    trg_hk.TriggerHookManagerClient
	logger *slog.Logger
}

func NewClient(ctx context.Context,
	logger *slog.Logger, addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "trigger_hook_grpc.NewClient"

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpc_retry.WithMax(uint(retriesCount)),
		grpc_retry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpc_logging.Option{
		grpc_logging.WithLogOnEvents(grpc_logging.PayloadReceived, grpc_logging.PayloadSent),
	}
	fmt.Println(addr, "sd")
	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpc_logging.UnaryClientInterceptor(InterceptorLogger(logger), logOpts...),
			grpc_retry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Client{
		api:    trg_hk.NewTriggerHookManagerClient(cc),
		logger: logger,
	}, nil
}

func InterceptorLogger(logger *slog.Logger) grpc_logging.Logger {
	return grpc_logging.LoggerFunc(func(ctx context.Context, lvl grpc_logging.Level, msg string, fields ...any) {
		logger.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) ScheduleDeletion(ctx context.Context, bannerID, tagID, featureID int64) (bool, error) {
	const op = "trigger_hook_grpc.ScheduleDeletion"
	resp, err := c.api.ScheduleDeletion(ctx, &trg_hk.CreateDeletionRequest{
		BannerID:  bannerID,
		TagID:     tagID,
		FeatureID: featureID,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %v", op, err)
	}

	return resp.Success, nil
}

//func (c *Client) SchedulePartialDeletion(ctx context.Context, tagID, featureID int64) (bool, error) {
//	const op = "grpc.ScheduleFullDeletion"
//
//	resp, err := c.api.SchedulePartialDeletion(ctx, &trg_hk.CreatePartialDeletionRequest{
//		TagID:     tagID,
//		FeatureID: featureID,
//	})
//	if err != nil {
//		return false, fmt.Errorf("%s: %v", op, err)
//	}
//
//	return resp.Success, nil
//}
