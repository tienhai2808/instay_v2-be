package initialization

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	rAddr := cfg.Host + fmt.Sprintf(":%d", cfg.Port)

	options := &redis.Options{
		Addr:     rAddr,
		Password: cfg.Password,
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	}
	if cfg.UseSSL {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
