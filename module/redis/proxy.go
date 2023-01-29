package redis

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/19
*@description:
 */
import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"gitlab.com/canyinxinxi/main-server/config"
	"gitlab.com/canyinxinxi/main-server/module/utils"
	"time"
)

var proxyRedis *redis.Client

const (
	RedisMaxUserIdKey = "max_user_id_key"
	RedisMaxJobIdKey  = "max_job_id_key"
)

func Init(sequence int) (err error) {
	/**
	redis 初始化
	*/
	//core - 核心数据缓存
	optionsCore := utils.Must(redis.ParseURL(config.GetConfig().Redis.Url)).(*redis.Options)

	//DB-0 - 核心数据缓存
	optionsCoreCore := *optionsCore
	optionsCoreCore.DB = 0

	optionsCoreCore.PoolSize = config.GetConfig().Redis.PoolSize
	proxyRedis = redis.NewClient(&optionsCoreCore)
	err = proxyRedis.Ping(context.Background()).Err()
	if err != nil {
		return err
	}

	return nil

}
func GetRedis() *redis.Client {
	return proxyRedis
}

func GetMaxUserId() int64 {
	return proxyRedis.Incr(context.Background(), RedisMaxUserIdKey).Val()
}
func GetMaxJobId() int64 {
	return proxyRedis.Incr(context.Background(), RedisMaxJobIdKey).Val()
}
func SetUserToken(userId int64) string {
	id := xid.New().String()
	proxyRedis.Set(context.Background(), id, userId, time.Hour*24*7)
	return id
}
func GetUserToken(token string) string {
	return proxyRedis.Get(context.Background(), token).Val()
}
