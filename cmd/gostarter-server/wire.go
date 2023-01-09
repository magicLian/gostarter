//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"patrol/cmd/patrol-server/backgroundSvc"
	"patrol/pkg/api"
	"patrol/pkg/middleware"
	asset_register "patrol/pkg/services/asset-register"
	"patrol/pkg/services/cronTask"
	"patrol/pkg/services/dashboard"
	"patrol/pkg/services/health"
	"patrol/pkg/services/org"
	"patrol/pkg/services/patrolRule"
	"patrol/pkg/services/patrolTask"
	"patrol/pkg/services/s3"
	"patrol/pkg/services/sqlstore"
	"patrol/pkg/services/sso"
	"patrol/pkg/services/systemRemindRule"
	"patrol/pkg/services/taskTemplate"
	"patrol/pkg/services/userManagement"
	"patrol/pkg/setting"
)

var wireSet = wire.NewSet(
	setting.ProvideSettingCfg,
	sqlstore.ProvideSqlStore,
	api.ProvideHttpServer,
	NewPatrolServer,
	middleware.ProvideMiddleWare,
	backgroundSvc.ProviceBackgroupServiceRegistry,
	cronTask.ProvideCronTasksService,
	health.ProvideHealthService,
	patrolRule.ProvicePatrolRuleService,
	patrolTask.ProvidePartolTaskService,
	taskTemplate.ProvideTaskTemplateService,
	s3.ProvideS3Service,
	sso.ProvideSsoService,
	systemRemindRule.ProvideSystemRemindRule,
	org.ProvideOrgService,
	asset_register.ProvideAssetRegisterInterface,
	userManagement.ProvideUserMgmtService,
	dashboard.ProvideDashboardService,
)

func InitGoStarterWire(cmd *setting.CommandLineArgs) (*PatrolServer, error) {
	wire.Build(wireSet)
	return &PatrolServer{}, nil
}
