package backgroundSvc

import (
	"context"

	"github.com/magicLian/gostarter/pkg/api"
	"github.com/magicLian/gostarter/pkg/services/cronTask"
)

type BackgroundService interface {
	Run(ctx context.Context) error
}

type BackgroundServiceRegistry struct {
	services []BackgroundService
}

func ProviceBackgroupServiceRegistry(httpServer *api.HTTPServer, cronTask *cronTask.CronTasksService) *BackgroundServiceRegistry {
	return NewBackgroundServiceRegistry(
		httpServer,
		cronTask,
	)
}

func NewBackgroundServiceRegistry(services ...BackgroundService) *BackgroundServiceRegistry {
	return &BackgroundServiceRegistry{
		services: services,
	}
}

func (r *BackgroundServiceRegistry) GetServices() []BackgroundService {
	return r.services
}
