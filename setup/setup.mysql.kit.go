package setup

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/signalfx/splunk-otel-go/instrumentation/github.com/jmoiron/sqlx/splunksqlx"

	kit "github.com/eden/go-biz-kit"
	"github.com/eden/go-biz-kit/config/def"
	"github.com/eden/go-biz-kit/database"
)

// NewMySQLDatabase 创建 MySQL 客户端
func NewMySQLDatabase(conf *def.Configuration, logger log.Logger) (kit.MySQL, error) {
	return NewSQLDatabase(conf, logger)
}

// NewSQLDatabase 创建满足 SQL 规范的客户端
func NewSQLDatabase(conf *def.Configuration, logger log.Logger) (kit.Database, error) {
	config := &conf.Database

	if !config.GetEnable() {
		return nil, nil
	}

	var err error = nil
	var db *sqlx.DB

	driver := config.GetDriver()
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	if conf.Tracing.GetEnable() {
		db, err = splunksqlx.Open(driver, config.GetAddr())
	} else {
		db, err = sqlx.Open(driver, config.GetAddr())
	}

	logHelper := log.NewHelper(log.With(logger, "module", driver))

	if err != nil {
		logHelper.Fatalw("msg", driver+" connect failed", "err", err)
	}

	if err1 := db.Ping(); err1 != nil {
		logHelper.Fatalw("msg", driver+" ping failed", "err", err1)
	}

	if config.GetConnMaxIdleTime() != nil {
		db.SetConnMaxIdleTime(config.GetConnMaxIdleTime().AsDuration())
	}

	if config.GetConnMaxLifetime() != nil {
		db.SetConnMaxLifetime(config.GetConnMaxLifetime().AsDuration())
	}

	db.SetMaxIdleConns(int(config.GetMaxPoolIdleSize()))
	db.SetMaxOpenConns(int(config.GetMaxConnection()))

	return database.NewDB(db), nil
}
