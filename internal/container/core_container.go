package container

import (
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/initialization"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/jwt"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/rabbitmq"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/redis"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/smtp"
)

func (c *Container) initCore() (err error) {
	c.Log, err = initialization.InitZap(c.cfg.Log)
	if err != nil {
		return err
	}

	c.DB, err = initialization.InitDatabase(c.cfg.PostgreSQL)
	if err != nil {
		return err
	}

	c.cache, err = initialization.InitRedis(c.cfg.Redis)
	if err != nil {
		return err
	}

	c.stor, err = initialization.InitS3(c.cfg.MinIO)
	if err != nil {
		return err
	}

	c.mq, err = initialization.InitRabbitMQ(c.cfg.RabbitMQ)
	if err != nil {
		return err
	}

	c.IDGen, err = initialization.InitSnowFlake()
	if err != nil {
		return err
	}

	c.jwtPro = jwt.NewJWTProvider(c.cfg.JWT)

	c.cachePro = redis.NewCacheProvider(c.cache)

	c.MQPro = rabbitmq.NewMessageQueueProvider(c.mq.Conn, c.mq.Chan, c.Log)

	c.SMTPPro = smtp.NewSMTPProvider(c.cfg.SMTPConfig)

	return nil
}
