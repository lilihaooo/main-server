package producet

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gitlab.com/canyinxinxi/main-server/config"
	redis2 "gitlab.com/canyinxinxi/main-server/module/redis"
	"gitlab.com/canyinxinxi/main-server/pb/index"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/25
*@description:
 */

//var Producet *redis.PubSub
//
//func InitProducetIndex(channel string) (err error) {
//	Producet = redis2.GetRedis().PSubscribe(context.Background(), channel)
//	return nil
//
//}

func SendUserRr(rr *index.User) {
	if rr == nil {
		log.Error().Msgf("user is null %v", rr)
		return
	}
	data, err := json.Marshal(rr)
	if err != nil {
		log.Error().Msgf("user:  %v proto err : %v", rr, err)
	} else {
		updateNotify := new(index.UpdateNotify)
		updateNotify.Transaction = &index.UpdateNotify_Transaction{
			Type:    index.UpdateNotify_User,
			Size:    0,
			Keys:    nil,
			Message: string(data),
		}
		sendData, errs := json.Marshal(updateNotify)
		if errs != nil {
			log.Error().Msgf("user:  %v proto err : %v", rr, err)
		} else {
			Send(string(sendData))
		}

	}
}
func SendJobRr(rr *index.Job) {
	if rr == nil {
		log.Error().Msgf("user is null %v", rr)
		return
	}
	data, err := json.Marshal(rr)
	if err != nil {
		log.Error().Msgf("user:  %v proto err : %v", rr, err)
	} else {
		updateNotify := new(index.UpdateNotify)
		updateNotify.Transaction = &index.UpdateNotify_Transaction{
			Type:    index.UpdateNotify_Job,
			Size:    0,
			Keys:    nil,
			Message: string(data),
		}
		sendData, errs := json.Marshal(updateNotify)
		if errs != nil {
			log.Error().Msgf("user:  %v proto err : %v", rr, err)
		} else {
			Send(string(sendData))
		}

	}
}
func Send(message string) {
	redis2.GetRedis().RPush(context.Background(),
		config.GetConfig().Redis.Producet.Inverted.Topic, message)
}
