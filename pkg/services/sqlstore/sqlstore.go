package sqlstore

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/magicLian/gostarter/pkg/setting"
	"github.com/magicLian/gostarter/pkg/util"

	glog "github.com/magicLian/gostarter/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	POSTGRES = "postgres"
	SQLITE   = "sqlite3"

	SqlSortDesc = "DESC"
	sqlSortAsc  = "ASC"
)

type DBSecretCfg struct {
	DBSchemaNmae string
	DBGroupName  string
}

type SqlStoreInterface interface {
	HealthStore
	GetDB() *gorm.DB
}

type SqlStore struct {
	cfg *setting.Cfg
	db  *gorm.DB
	SqlStoreInterface
	dbCfg DatabaseConfig
	log   glog.Logger
}

type CommonSqlstore struct {
	db  *gorm.DB
	log log.Logger
}

type DatabaseConfig struct {
	Type             string
	ConnectionString string
	User             string
	Password         string
	Host             string
	Name             string
	SchemaName       string
	GroupName        string
	MaxOpenConn      int
	MaxIdleConn      int
	ConnMaxLifetime  int
	LogQueries       bool
	UrlQueryParams   map[string][]string
	SslMode          string
	CaCertPath       string
	ClientKeyPath    string
	ClientCertPath   string

	//sqlite3
	Path      string
	CacheMode string

	//table perfix
	TableNamePerfix string
}

func ProvideSqlStore(cfg *setting.Cfg) (SqlStoreInterface, error) {
	ss := &SqlStore{
		cfg: cfg,
		log: glog.New("sqlstore"),
	}
	ss.readConfig()

	db, err := ss.getEngine()
	if err != nil {
		return nil, fmt.Errorf("Fail to connect to database: %v", err)
	}
	ss.db = db

	//create table
	// if err = ss.createTable(models.ResourceChangeRecord{}); err != nil {
	// 	return nil, err
	// }

	ss.SqlStoreInterface = NewCommonSqlstore(db)

	return ss, nil
}

func NewCommonSqlstore(db *gorm.DB) SqlStoreInterface {
	return &SqlStore{
		db:  db,
		log: glog.New("common sqlstore"),
	}
}

func (ss *SqlStore) GetDB() *gorm.DB {
	return ss.db
}

func (ss *SqlStore) getEngine() (*gorm.DB, error) {
	connStr, err := ss.buildConnectionString()
	if err != nil {
		return nil, err
	}
	ss.log.Infof("Connecting to DB", "dbtype", ss.dbCfg.Type)
	config := &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   ss.dbCfg.TableNamePerfix,
		SingularTable: true,
	}}
	if ss.dbCfg.LogQueries {
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,
			},
		)
	}
	engine, err := gorm.Open(postgres.Open(connStr), config)
	if err != nil {
		ss.log.Errorf("Failed to connection to db, err:[%s]", err.Error())
		return nil, err
	}

	if ss.dbCfg.SchemaName != "" {
		//ALTER ROLE \"" + ss.dbCfg.User + "\" SET search_path TO \"" + ss.dbCfg.SchemaName + "\";
		err = engine.Exec("CREATE SCHEMA IF NOT EXISTS \"" + ss.dbCfg.SchemaName + "\" AUTHORIZATION " + ss.dbCfg.GroupName + ";").Error
		if err != nil {
			ss.log.Errorf("Failed to create schema.", "schema", ss.dbCfg.SchemaName, "group", ss.dbCfg.GroupName, "error", err.Error())
			return nil, err
		}
	}

	db, err := engine.DB()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Second * time.Duration(ss.dbCfg.ConnMaxLifetime))
	db.SetMaxOpenConns(ss.dbCfg.MaxOpenConn)
	db.SetMaxIdleConns(ss.dbCfg.MaxIdleConn)

	return engine, nil
}

func (ss *SqlStore) buildConnectionString() (string, error) {
	cnnstr := ss.dbCfg.ConnectionString
	if cnnstr != "" {
		return ss.dbCfg.ConnectionString, nil
	}
	switch ss.dbCfg.Type {
	case POSTGRES:
		addr, err := util.SplitHostPortDefault(ss.dbCfg.Host, "127.0.0.1", "5432")
		if err != nil {
			return "", fmt.Errorf("Invalid host specifier '%s',error %s", ss.dbCfg.Host, err.Error())
		}
		if ss.dbCfg.Password == "" {
			ss.dbCfg.Password = "''"
		}
		if ss.dbCfg.User == "" {
			ss.dbCfg.User = "''"
		}
		cnnstr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s ", ss.dbCfg.User, ss.dbCfg.Password, addr.Host, addr.Port, ss.dbCfg.Name, ss.dbCfg.SslMode)
		if ss.dbCfg.SslMode != "disable" {
			cnnstr += fmt.Sprintf(" sslcert=%s sslkey=%s sslrootcert=%s ", ss.dbCfg.ClientCertPath, ss.dbCfg.ClientKeyPath, ss.dbCfg.CaCertPath)
		}
		if ss.dbCfg.SchemaName != "" {
			cnnstr += " search_path=" + ss.dbCfg.SchemaName
		}
	case SQLITE:
		if !filepath.IsAbs(ss.dbCfg.Path) {
			ss.dbCfg.Path = filepath.Join(ss.cfg.DataPath, ss.dbCfg.Path)
		}
		if err := os.MkdirAll(path.Dir(ss.dbCfg.Path), os.ModePerm); err != nil {
			return "", err
		}

		cnnstr = fmt.Sprintf("file:%s?cache=%s&mode=rwc", ss.dbCfg.Path, ss.dbCfg.CacheMode)

		cnnstr += ss.buildExtraConnectionString('&')
	default:
		return "", fmt.Errorf("Unknown database type: %s", ss.dbCfg.Type)
	}
	return cnnstr, nil
}

func (ss *SqlStore) readConfig() {
	ss.readLocalConfig()

	ss.dbCfg.SchemaName = util.GetEnvOrIniValue(ss.cfg.Raw, "database", "schema_name")
	ss.dbCfg.GroupName = util.GetEnvOrIniValue(ss.cfg.Raw, "database", "group_name")
	ss.dbCfg.MaxOpenConn = util.SetDefaultInt(util.GetEnvOrIniValue(ss.cfg.Raw, "database", "max_open_conn"), 100)
	ss.dbCfg.MaxIdleConn = util.SetDefaultInt(util.GetEnvOrIniValue(ss.cfg.Raw, "database", "max_idle_conn"), 100)
	ss.dbCfg.ConnMaxLifetime = util.SetDefaultInt(util.GetEnvOrIniValue(ss.cfg.Raw, "database", "conn_max_life_time"), 14400)
	ss.dbCfg.LogQueries = util.SetDefaultBool(util.GetEnvOrIniValue(ss.cfg.Raw, "database", "log_queries"), false)
	ss.dbCfg.SslMode = util.SetDefaultString(util.GetEnvOrIniValue(ss.cfg.Raw, "database", "ssl_mode"), "disable")
	ss.dbCfg.CaCertPath = util.GetEnvOrIniValue(ss.cfg.Raw, "database", "ca_cert_path")
	ss.dbCfg.ClientKeyPath = util.GetEnvOrIniValue(ss.cfg.Raw, "database", "client_key_path")
	ss.dbCfg.ClientCertPath = util.GetEnvOrIniValue(ss.cfg.Raw, "database", "client_cert_path")
	ss.dbCfg.Path = util.SetDefaultString(util.GetEnvOrIniValue(ss.cfg.Raw, "database", "path"), "data/patrol.db")
	ss.dbCfg.CacheMode = util.SetDefaultString(util.GetEnvOrIniValue(ss.cfg.Raw, "database", "cache_mode"), "private")
	ss.dbCfg.TableNamePerfix = util.SetDefaultString(util.GetEnvOrIniValue(ss.cfg.Raw, "database", "table_perfix"), "p_")
}

func (ss *SqlStore) readLocalConfig() {
	dbUrl := util.GetEnvOrIniValue(ss.cfg.Raw, "database", "url")
	if dbUrl != "" {
		ss.dbCfg.ConnectionString = dbUrl
		return
	}
	ss.dbCfg.Type = util.SetDefaultString(util.GetEnvOrIniValue(ss.cfg.Raw, "database", "type"), "postgres")
	ss.dbCfg.Host = util.GetEnvOrIniValue(ss.cfg.Raw, "database", "host")
	ss.dbCfg.Name = util.GetEnvOrIniValue(ss.cfg.Raw, "database", "name")
	ss.dbCfg.User = util.GetEnvOrIniValue(ss.cfg.Raw, "database", "user")
	ss.dbCfg.Password = util.GetEnvOrIniValue(ss.cfg.Raw, "database", "password")
}

func (ss *SqlStore) buildExtraConnectionString(sep rune) string {
	if ss.dbCfg.UrlQueryParams == nil {
		return ""
	}

	var sb strings.Builder
	for key, values := range ss.dbCfg.UrlQueryParams {
		for _, value := range values {
			sb.WriteRune(sep)
			sb.WriteString(key)
			sb.WriteRune('=')
			sb.WriteString(value)
		}
	}
	return sb.String()
}

func (ss *SqlStore) createTable(dst interface{}) error {
	return ss.db.Transaction(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(&dst)
		if err != nil {
			return err
		}
		if ss.dbCfg.GroupName != "" {
			structName := reflect.ValueOf(dst).Type().Name()
			sql := "ALTER TABLE " + ss.GetFullTableName(strings.ToLower(util.Camel2Case(structName))) + " OWNER TO " + ss.dbCfg.GroupName
			err = tx.Exec(sql).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (db *SqlStore) GetFullTableName(tablename string) string {
	if db.dbCfg.TableNamePerfix != "" {
		return fmt.Sprintf("%s%s", db.dbCfg.TableNamePerfix, tablename)
	}
	return tablename
}
