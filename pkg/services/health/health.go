package health

import (
	"github.com/magicLian/gostarter/pkg/services/sqlstore"
	"github.com/magicLian/gostarter/pkg/setting"

	"github.com/magicLian/gostarter/pkg/models"
)

type Health interface {
	HealthCheck() (bool, *models.Health)
}

func ProvideHealthService(storeInterface sqlstore.SqlStoreInterface) Health {
	return &HealthService{
		sqlstore: storeInterface,
	}
}

type HealthService struct {
	sqlstore sqlstore.SqlStoreInterface
}

func (h *HealthService) HealthCheck() (bool, *models.Health) {
	allOk := true
	result := &models.Health{}
	if err := h.sqlstore.DBPing(); err != nil {
		result.Database = models.DATABASE_FAILING
		allOk = false
	} else {
		result.Database = models.DATABASE_OK
	}
	result.ApiVersion = setting.APIVersion
	return allOk, result
}
