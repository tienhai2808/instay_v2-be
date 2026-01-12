package container

import "github.com/InstayPMS/backend/internal/infrastructure/api/http/handler"

func (c *Container) initHandlers() {
	c.FileHdl = handler.NewFileHandler(c.fileUC)

	c.AuthHdl = handler.NewAuthHandler(c.cfg, c.authUC)

	c.UserHdl = handler.NewUserHandler(c.userUC)

	c.DepartmentHdl = handler.NewDepartmentHandler(c.departmentUC)
}
