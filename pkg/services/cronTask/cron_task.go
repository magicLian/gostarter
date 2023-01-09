package cronTask

import (
	"context"
	"os"
	"path"

	"github.com/magicLian/gostarter/pkg/setting"
	"github.com/magicLian/gostarter/pkg/util"

	"github.com/magicLian/gostarter/pkg/log"
)

type CronTasksService struct {
	log      log.Logger
	cfg      *setting.Cfg
	timeZone string
	//add cron string
}

func ProvideCronTasksService(cfg *setting.Cfg) *CronTasksService {
	service := &CronTasksService{
		log:      log.New("cron task"),
		cfg:      cfg,
		timeZone: util.SetDefaultString(util.GetEnvOrIniValue(cfg.Raw, "server", "time_zone"), "Local"),
	}

	return service
}

func (c *CronTasksService) Run(ctx context.Context) error {
	c.log.Infof("init cron tasks...")
	c.setTimeZoneEnv()
	/*location, err := util.GetTimeLocation(c.timeZone)
	if err != nil {
		return err
	}

	parser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	cronTasks := cron.New(cron.WithParser(parser), cron.WithLocation(location))

	cronTasks.Start()*/
	return nil
}

func (c *CronTasksService) setTimeZoneEnv() {
	if err := os.Setenv("ZONEINFO", path.Join(setting.HomePath, "public/tzdata/data.zip")); err != nil {
		c.log.Errorf("init os timezone file failed")
	}
}
