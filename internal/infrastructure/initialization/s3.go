package initialization

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitS3(cfg config.MinIOConfig) (*s3.Client, error) {
	staticCredentials := credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		"",
	)

	aCfg, err := awsCfg.LoadDefaultConfig(
		context.TODO(),
		awsCfg.WithRegion(cfg.Region),
		awsCfg.WithCredentialsProvider(staticCredentials),
		awsCfg.WithRetryMaxAttempts(3),
		awsCfg.WithRetryMode(aws.RetryModeStandard),
		awsCfg.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	protocol := "http"
	if cfg.UseSSL {
		protocol = protocol + "s"
	}
	client := s3.NewFromConfig(aCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(fmt.Sprintf("%s://%s", protocol, cfg.Endpoint))
	})

	return client, nil
}
