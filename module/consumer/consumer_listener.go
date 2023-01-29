package consumer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"gitlab.com/canyinxinxi/main-server/db"
	redis2 "gitlab.com/canyinxinxi/main-server/module/redis"
	"gitlab.com/canyinxinxi/main-server/pb/index"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/20
*@description:
 */
var Consumer *redis.PubSub

func InitConsumerIndex(channel string) (err error) {
	Consumer = redis2.GetRedis().Subscribe(context.Background(), channel)
	go listen()
	return nil

}

func listen() {
	for {
		message, err := Consumer.Receive(context.Background())
		if err != nil {
			continue
		}
		// 检测收到的消息类型
		switch message.(type) {
		case *redis.Subscription:
			// 订阅成功
		case *redis.Message:
			// 处理收到的消息
			// 这里需要做一下类型转换
			data := message.(*redis.Message)
			updateNotify := new(index.UpdateNotify)

			if err := proto.Unmarshal([]byte(data.Pattern), updateNotify); err != nil {
				log.Error().Msgf("增量同步异常:%s", err.Error())
			}

			if updateNotify.Transaction == nil {
				log.Error().Msgf("topic %s transaction is nil", message)
				return
			}

			if len(updateNotify.Transaction.Keys) == 0 {
				log.Error().Msgf("keys len is 0")
				return
			}

			db.Update(updateNotify)
			// 打印收到的消息
		case *redis.Pong:
			// 收到Pong消息
		default:
			continue

		}
	}
}
