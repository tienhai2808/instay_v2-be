package container

import (
	authUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/auth"
	departmentUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/department"
	fileUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/file"
	userUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/user"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm"
)

func (c *Container) initLogic() {
	c.UserRepo = orm.NewUserRepository(c.DB.Gorm)
	c.tokenRepo = orm.NewTokenRepository(c.DB.Gorm)
	c.departmentRepo = orm.NewDepartmentRepository(c.DB.Gorm)

	c.fileUC = fileUC.NewFileUseCase(c.cfg, c.stor, c.Log)
	c.authUC = authUC.NewAuthUseCase(c.cfg.JWT, c.DB.Gorm, c.Log, c.IDGen, c.jwtPro, c.cachePro, c.MQPro, c.UserRepo, c.tokenRepo)
	c.userUC = userUC.NewUserUseCase(c.DB.Gorm, c.Log, c.IDGen, c.cachePro, c.UserRepo, c.departmentRepo, c.tokenRepo)
	c.departmentUC = departmentUC.NewDepartmentUseCase(c.Log, c.IDGen, c.departmentRepo)
}
