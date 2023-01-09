package setting

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/magicLian/gostarter/pkg/util"

	"gitlab.wise-paas.com/ifactory/patrol-models/log"
	"gopkg.in/ini.v1"
)

type Scheme string

const (
	DEV      = "development"
	PROD     = "production"
	APP_NAME = "patrol-api"
)

var (
	Env             = DEV
	HomePath        string
	configFiles     []string
	AppUrl          string
	HttpPort        string
	Domain          string
	ApplicationName string

	APIVersion string

	SsoApiUrl      string
	SubscriptionId string

	AssetRegisterApiUrl string
)

type Cfg struct {
	Raw *ini.File

	//ViperCfg *viper.Viper
	DataPath string
	LogsPath string
}

type CommandLineArgs struct {
	Config   string
	HomePath string
	Args     []string
}

func makeAbsolute(path string, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func evalEnvVarExpression(value string) string {
	regex := regexp.MustCompile(`\${(\w+)}`)
	return regex.ReplaceAllStringFunc(value, func(envVar string) string {
		envVar = strings.TrimPrefix(envVar, "${")
		envVar = strings.TrimSuffix(envVar, "}")
		envValue := os.Getenv(envVar)

		// if env variable is hostname and it is empty use os.Hostname as default
		if envVar == "HOSTNAME" && envValue == "" {
			envValue, _ = os.Hostname()
		}

		return envValue
	})
}

func evalConfigValues(file *ini.File) {
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			key.SetValue(evalEnvVarExpression(key.Value()))
		}
	}
}

func (cfg *Cfg) LoadConfiguration(args *CommandLineArgs) (*ini.File, error) {
	var err error
	// load config defaults
	defaultConfigFile := path.Join(HomePath, "conf/defaults.ini")
	configFiles = append(configFiles, defaultConfigFile)

	// check if config file exists
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		os.Exit(1)
	}
	parsedFile, err := ini.Load(defaultConfigFile)
	if err != nil {
		fmt.Printf("Failed to parse defaults.ini, %v\n", err)
		os.Exit(1)
		return nil, err
	}

	parsedFile.BlockMode = false

	// evaluate config values containing environment variables
	evalConfigValues(parsedFile)

	cfg.initLogging(parsedFile)

	cfg.Raw = parsedFile

	return parsedFile, err
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func SetHomePath(args *CommandLineArgs) {

	if args.HomePath != "" {
		HomePath = args.HomePath
		return
	}

	HomePath, _ = filepath.Abs(".")
	//	// check if homepath is correct
	if pathExists(filepath.Join(HomePath, "conf/defaults.ini")) {
		return
	}

	// try down one path
	if pathExists(filepath.Join(HomePath, "../conf/defaults.ini")) {
		HomePath = filepath.Join(HomePath, "../")
	}
}

func NewCfg() *Cfg {
	return &Cfg{
		Raw: ini.Empty(),
	}
}

func ProvideSettingCfg(args *CommandLineArgs) (*Cfg, error) {
	cfg := NewCfg()
	if err := cfg.Load(args); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Cfg) Load(args *CommandLineArgs) error {
	SetHomePath(args)

	_, err := cfg.LoadConfiguration(args)
	if err != nil {
		return err
	}
	/*	if err := cfg.LoadViperConfig(); err != nil {
		return err
	}*/
	Env = util.SetDefaultString(util.GetEnvOrIniValue(cfg.Raw, "app_mode", "development"), DEV)
	HttpPort = util.SetDefaultString(util.GetEnvOrIniValue(cfg.Raw, "server", "http_port"), "8080")
	ApplicationName = APP_NAME

	APIVersion = util.SetDefaultString(util.GetEnvOrIniValue(cfg.Raw, "server", "apiversion"), "0.0.1")

	SsoApiUrl = util.GetEnvOrIniValue(cfg.Raw, "sso", "api_url")
	if strings.HasSuffix(SsoApiUrl, "/v4.0") {
		SsoApiUrl = strings.Trim(SsoApiUrl, "/v4.0")
	}

	/*SubscriptionId = cfg.ViperCfg.GetString("sso.subscription_id")
	Domain = cfg.ViperCfg.Get("sso.domain")*/
	SubscriptionId = util.GetEnvOrIniValue(cfg.Raw, "sso", "subscription_id")
	Domain = util.GetEnvOrIniValue(cfg.Raw, "sso", "domain")
	if strings.HasPrefix(Domain, ".") {
		Domain = Domain[1:]
	}
	AppUrl = util.GetEnvOrIniValue(cfg.Raw, "server", "app_url")
	if AppUrl == "" {
		AppUrl = getDefaultAppUrl()
	}

	AssetRegisterApiUrl = util.GetEnvOrIniValue(cfg.Raw, "asset-register", "api_url")
	return nil
}

func getDefaultAppUrl() string {
	namespaceEnv := os.Getenv("namespace")
	clusterEnv := os.Getenv("cluster")
	if namespaceEnv != "" && clusterEnv != "" {
		return fmt.Sprintf("%s-%s-%s.%s", ApplicationName, namespaceEnv, clusterEnv, Domain)
	}
	return fmt.Sprintf("localhost:%s", HttpPort)
}

/*func (cfg *Cfg) LoadViperConfig() error {
	configFile := path.Join(HomePath, "conf/defaults.ini")
	v := viper.New()
	v.SetConfigFile(configFile)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	v.AutomaticEnv()
	cfg.ViperCfg = v
	return nil
}
*/
func (cfg *Cfg) initLogging(file *ini.File) {
	logModes := strings.Split(util.SetDefaultString(util.GetEnvOrIniValue(cfg.Raw, "log", "mode"), "console"), ",")
	if len(logModes) == 1 {
		logModes = strings.Split(util.SetDefaultString(util.GetEnvOrIniValue(cfg.Raw, "log", "mode"), "console"), " ")
	}

	cfg.LogsPath = makeAbsolute(util.SetDefaultString(util.GetEnvOrIniValue(cfg.Raw, "paths", "logs"), ""), HomePath)
	logLevel := util.SetDefaultString(util.GetEnvOrIniValue(file, "log", "level"), "info")
	log.ReadLoggingConfig(logModes, logLevel, cfg.LogsPath, file)
}
