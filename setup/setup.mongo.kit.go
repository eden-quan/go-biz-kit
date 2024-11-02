package setup

import (
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/mongo"

	mongopkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/mongo"

	kit "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git"
	config "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/config/def"
)

type mongoDBImpl struct {
	client *mongo.Client
	db     *mongo.Database
}

func newMongoDB(db *mongo.Database, client *mongo.Client) kit.MongoDB {
	return &mongoDBImpl{
		client: client,
		db:     db,
	}
}

func (db *mongoDBImpl) Get() *mongo.Database {
	return db.db
}
func (db *mongoDBImpl) GetDB(name string) *mongo.Database {
	return db.client.Database(name)
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
	return newMongoDB(db, client), nil
}
