package setup

import (
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/mongo"

	mongopkg "github.com/eden/go-kratos-pkg/mongo"

	kit "github.com/eden/go-biz-kit"
	config "github.com/eden/go-biz-kit/config/def"
)

type mongoDBImpl struct {
	db *mongo.Database
}

func newMongoDB(db *mongo.Database) kit.MongoDB {
	return &mongoDBImpl{
		db: db,
	}
}

func (db *mongoDBImpl) Get() *mongo.Database {
	return db.db
}

// NewMongoDB 创建 MongoDB 客户端 mongo database
func NewMongoDB(conf *config.Configuration, logger log.Logger) (kit.MongoDB, error) {

	if !conf.Mongo.GetEnable() {
		return nil, nil
	}

	mongoConfig := &conf.Mongo

	c := &mongopkg.Config{
		Addr:              mongoConfig.GetAddress(),
		MaxPoolSize:       mongoConfig.GetMaxPoolSize(),
		MinPoolSize:       mongoConfig.GetMinPoolSize(),
		MaxConnecting:     mongoConfig.GetMaxConnection(),
		ConnectTimeout:    mongoConfig.GetConnectTimeout(),
		HeartbeatInterval: mongoConfig.GetHeartbeatInterval(),
		MaxConnIdleTime:   mongoConfig.GetMaxConnIdleTime(),
		Timeout:           mongoConfig.GetTimeout(),
		Hosts:             mongoConfig.GetHosts(),
		Debug:             mongoConfig.GetDebug(),
	}
	client := mongopkg.NewMongoClient(c, logger)

	db := client.Database(mongoConfig.Database)
	return newMongoDB(db), nil
}
