package container

import (
	authUC "github.com/InstayPMS/backend/internal/application/usecase/auth"
	departmentUC "github.com/InstayPMS/backend/internal/application/usecase/department"
	fileUC "github.com/InstayPMS/backend/internal/application/usecase/file"
	userUC "github.com/InstayPMS/backend/internal/application/usecase/user"
)

func (c *Container) initUseCases() {
	c.fileUC = fileUC.NewFileUseCase(c.cfg, c.stor, c.Log)

	c.authUC = authUC.NewAuthUseCase(c.cfg.JWT, c.db.Gorm, c.Log, c.idGen, c.jwtPro, c.cachePro, c.MQPro, c.userRepo, c.tokenRepo)

	c.userUC = userUC.NewUserUseCase(c.Log, c.idGen, c.userRepo, c.departmentRepo)

	c.departmentUC = departmentUC.NewDepartmentUseCase(c.Log, c.idGen, c.departmentRepo)
}
