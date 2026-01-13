package container

import (
	"log"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	authUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/auth"
	departmentUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/department"
	fileUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/file"
	userUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/user"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	httpHdl "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	httpMid "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/initialization"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/rabbitmq"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/smtp"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake/v2"
	"go.uber.org/zap"
)

type Container struct {
	cfg               *config.Config
	Log               *zap.Logger
	DB                *initialization.Database
	cache             *redis.Client
	mq                *initialization.MQ
	stor              *minio.Client
	IDGen             *sonyflake.Sonyflake
	jwtPro            port.JWTProvider
	MQPro             port.MessageQueueProvider
	cachePro          port.CacheProvider
	SMTPPro           port.SMTPProvider
	UserRepo          repository.UserRepository
	tokenRepo         repository.TokenRepository
	departmentRepo    repository.DepartmentRepository
	fileUC            fileUC.FileUseCase
	authUC            authUC.AuthUseCase
	userUC            userUC.UserUseCase
	departmentUC      departmentUC.DepartmentUseCase
	FileHTTPHdl       *httpHdl.FileHandler
	AuthHTTPHdl       *httpHdl.AuthHandler
	UserHTTPHdl       *httpHdl.UserHandler
	DepartmentHTTPHdl *httpHdl.DepartmentHandler
	CtxHTTPMid        *httpMid.ContextMiddleware
	AuthHTTPMid       *httpMid.AuthMiddleware
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		cfg: cfg,
	}
}

func (c *Container) InitServer() error {
	if err := c.initCore(); err != nil {
		return err
	}

	c.initLogic()
	c.initAPI()

	return nil
}

func (c *Container) InitSeed() error {
	var err error
	c.Log, err = initialization.InitZap(c.cfg.Log)
	if err != nil {
		return err
	}

	c.DB, err = initialization.InitDatabase(c.cfg.PostgreSQL)
	if err != nil {
		return err
	}

	c.IDGen, err = initialization.InitSnowFlake()
	if err != nil {
		return err
	}

	c.UserRepo = orm.NewUserRepository(c.DB.Gorm)

	return nil
}

func (c *Container) InitConsumer() error {
	var err error
	c.Log, err = initialization.InitZap(c.cfg.Log)
	if err != nil {
		return err
	}

	c.mq, err = initialization.InitRabbitMQ(c.cfg.RabbitMQ)
	if err != nil {
		return err
	}

	c.MQPro = rabbitmq.NewMessageQueueProvider(c.mq.Conn, c.mq.Chan, c.Log)
	c.SMTPPro = smtp.NewSMTPProvider(c.cfg.SMTPConfig)

	return nil
}

func (c *Container) Cleanup() {
	if c.DB != nil {
		c.DB.Close()
	}
	if c.mq != nil {
		c.mq.Close()
	}

	log.Println("Container cleaned successfully")
}
