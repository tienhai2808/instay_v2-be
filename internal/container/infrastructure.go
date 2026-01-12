package container

import (
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/middleware"
	"github.com/InstayPMS/backend/internal/infrastructure/initialization"
	"github.com/InstayPMS/backend/internal/infrastructure/persistence/orm"
	"github.com/InstayPMS/backend/internal/infrastructure/provider/jwt"
	"github.com/InstayPMS/backend/internal/infrastructure/provider/rabbitmq"
	"github.com/InstayPMS/backend/internal/infrastructure/provider/redis"
	"github.com/InstayPMS/backend/internal/infrastructure/provider/smtp"
)

func (c *Container) initInfrastructure() error {
	log, err := initialization.InitZap(c.cfg.Log)
	if err != nil {
		return err
	}
	c.Log = log

	db, err := initialization.InitDatabase(c.cfg.PostgreSQL)
	if err != nil {
		return err
	}
	c.db = db

	cache, err := initialization.InitRedis(c.cfg.Redis)
	if err != nil {
		return err
	}
	c.cache = cache

	stor, err := initialization.InitMinIO(c.cfg.MinIO)
	if err != nil {
		return err
	}
	c.stor = stor

	mq, err := initialization.InitRabbitMQ(c.cfg.RabbitMQ)
	if err != nil {
		return err
	}
	c.mq = mq

	idGen, err := initialization.InitSnowFlake()
	if err != nil {
		return err
	}
	c.idGen = idGen

	c.jwtPro = jwt.NewJWTProvider(c.cfg.JWT)

	c.cachePro = redis.NewCacheProvider(c.cache)

	c.MQPro = rabbitmq.NewMessageQueueProvider(c.mq.Conn, c.mq.Chan, c.Log)

	c.SMTPPro = smtp.NewSMTPProvider(c.cfg.SMTPConfig)

	c.userRepo = orm.NewUserRepository(c.db.Gorm)

	c.tokenRepo = orm.NewTokenRepository(c.db.Gorm)

	c.departmentRepo = orm.NewDepartmentRepository(c.db.Gorm)

	c.CtxMid = middleware.NewContextMiddleware(log)

	c.AuthMid = middleware.NewAuthMiddleware(c.cfg.JWT, c.Log, c.jwtPro, c.cachePro)

	return nil
}
