package di

import fileUC "github.com/InstayPMS/backend/internal/application/usecase/file"

func (c *Container) initUseCases() {
	c.FileUseCase = fileUC.NewFileUseCase(c.Config, c.Storage, c.Log)
}