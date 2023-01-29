package resume

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

var resumeDaoInstance = &resumeDao{
	Data: sync.Map{},
}

func GetResumeDao() IResume {
	return resumeDaoInstance
}

type IResume interface {
	Set(info *index.Resume)
	Get(id string) (*index.Resume, bool)
	GetAll() []*index.Resume
	Del(id string)
	Reload(keys []string)
}

type resumeDao struct {
	Data sync.Map
}

func (dao *resumeDao) GetAll() []*index.Resume {
	ads := make([]*index.Resume, 0)
	dao.Data.Range(func(key, value interface{}) bool {
		ads = append(ads, value.(*index.Resume))
		return true
	})
	return ads
}

func (dao *resumeDao) Set(info *index.Resume) {
	if info.GetStatus() != 3 {
		dao.Data.Store(fmt.Sprint(info.GetId()), info)
	} else {
		dao.Data.Delete(fmt.Sprint(info.GetId()))
	}
}

func (dao *resumeDao) Get(id string) (*index.Resume, bool) {
	if info, ok := dao.Data.Load(id); !ok {
		return nil, ok
	} else {
		return info.(*index.Resume), ok
	}
}

func (dao *resumeDao) Del(id string) {
	dao.Data.Delete(id)
}

func (dao *resumeDao) Reload(keys []string) {
	for _, key := range keys {
		v, _ := redis.GetRedis().Get(context.Background(), key).Bytes()
		var resumeInfo index.Resume
		if err := proto.Unmarshal(v, &resumeInfo); err != nil {
			log.Error().Msgf("proto.Unmarshal fail. key: %v, %v", key, err.Error())
			continue
		}
		dao.Set(&resumeInfo)
	}
	log.Info().Msgf("resume更新成功. %v", keys)
}
