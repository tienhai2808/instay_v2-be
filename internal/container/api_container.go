package container

import (
	httpHdl "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	httpMid "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
)

func (c *Container) initAPI() {
	c.FileHTTPHdl = httpHdl.NewFileHandler(c.fileUC)
	c.AuthHTTPHdl = httpHdl.NewAuthHandler(c.cfg, c.authUC)
	c.UserHTTPHdl = httpHdl.NewUserHandler(c.userUC)
	c.DepartmentHTTPHdl = httpHdl.NewDepartmentHandler(c.departmentUC)

	c.CtxHTTPMid = httpMid.NewContextMiddleware(c.Log)
	c.AuthHTTPMid = httpMid.NewAuthMiddleware(c.cfg.JWT, c.Log, c.jwtPro, c.cachePro)
}
