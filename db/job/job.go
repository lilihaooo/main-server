package job

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

var jobDaoInstance = &jobDao{
	Data: sync.Map{},
}

const (
	MY_JOB = "my_job_%s"
)

func GetJobDao() IJob {
	return jobDaoInstance
}

type IJob interface {
	Set(info *index.Job)
	Get(id string) (*index.Job, bool)
	GetMy(id string) []string
	GetAll() []*index.Job
	Del(id string)
	Reload(keys []string)
}

type jobDao struct {
	Data sync.Map
}

func (dao *jobDao) GetAll() []*index.Job {
	ads := make([]*index.Job, 0)
	dao.Data.Range(func(key, value interface{}) bool {
		ads = append(ads, value.(*index.Job))
		return true
	})
	return ads
}

func (dao *jobDao) Set(info *index.Job) {
	if info.GetStatus() != 3 {
		dao.Data.Store(fmt.Sprint(info.GetId()), info)
	} else {
		dao.Data.Delete(fmt.Sprint(info.GetId()))
	}
}

func (dao *jobDao) Get(id string) (*index.Job, bool) {
	if info, ok := dao.Data.Load(id); !ok {
		return nil, ok
	} else {
		return info.(*index.Job), ok
	}
}

func (dao *jobDao) GetMy(userId string) []string {
	return redis.GetRedis().SMembers(context.Background(), fmt.Sprintf(MY_JOB, userId)).Val()
}

func (dao *jobDao) Del(id string) {
	dao.Data.Delete(id)
}

func (dao *jobDao) Reload(keys []string) {
	for _, key := range keys {
		v, _ := redis.GetRedis().Get(context.Background(), key).Bytes()
		var jobInfo index.Job
		if err := proto.Unmarshal(v, &jobInfo); err != nil {
			log.Error().Msgf("proto.Unmarshal fail. key: %v, %v", key, err.Error())
			continue
		}
		dao.Set(&jobInfo)
	}
	log.Info().Msgf("job更新成功. %v", keys)
}
