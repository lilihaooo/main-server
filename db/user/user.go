package user

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/20
*@description:
 */

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"gitlab.com/canyinxinxi/main-server/module/redis"
	"gitlab.com/canyinxinxi/main-server/pb/index"
	"sync"
)

var userDaoInstance = &userDao{
	Data: sync.Map{},
}

func GetUserDao() IUser {
	return userDaoInstance
}

type IUser interface {
	Set(info *index.User)
	Get(id string) (*index.User, bool)
	GetAll() []*index.User
	Del(id string)
	Reload(keys []string)
}

type userDao struct {
	Data sync.Map
}

func (dao *userDao) GetAll() []*index.User {
	ads := make([]*index.User, 0)
	dao.Data.Range(func(key, value interface{}) bool {
		ads = append(ads, value.(*index.User))
		return true
	})
	return ads
}

func (dao *userDao) Set(info *index.User) {
	dao.Data.Store(fmt.Sprint(info.GetPhone()), info)
	dao.Data.Store(fmt.Sprint(info.GetWechartOpenId()), info)
}

func (dao *userDao) Get(id string) (*index.User, bool) {
	if info, ok := dao.Data.Load(id); !ok {
		return nil, ok
	} else {
		return info.(*index.User), ok
	}
}

func (dao *userDao) Del(id string) {
	dao.Data.Delete(id)
}

func (dao *userDao) Reload(keys []string) {
	for _, key := range keys {
		v, _ := redis.GetRedis().Get(context.Background(), key).Bytes()
		var userInfo index.User
		if err := proto.Unmarshal(v, &userInfo); err != nil {
			log.Error().Msgf("proto.Unmarshal fail. key: %v, %v", key, err.Error())
			continue
		}
		dao.Set(&userInfo)
	}
	log.Info().Msgf("用户更新成功. %v", keys)
}
