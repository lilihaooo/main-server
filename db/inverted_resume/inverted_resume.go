package inverted_resume

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/20
*@description:
 */

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"gitlab.com/canyinxinxi/main-server/module/redis"
	"gitlab.com/canyinxinxi/main-server/pb/index"
	"strings"
	"sync"
)

var invertedResumeDaoInstance = &invertedResumeDao{
	Data: sync.Map{},
}

func GetInvertedResumeDao() IInvertedResume {
	return invertedResumeDaoInstance
}

type IInvertedResume interface {
	Set(invertedKey string, inverted []string)
	Get(id string) ([]string, bool)
	Del(id string)
	Reload(key string)
}

type invertedResumeDao struct {
	Data sync.Map
}

func (dao *invertedResumeDao) GetAll() []*index.Job {
	ads := make([]*index.Job, 0)
	dao.Data.Range(func(key, value interface{}) bool {
		ads = append(ads, value.(*index.Job))
		return true
	})
	return ads
}

func (dao *invertedResumeDao) Set(invertedKey string, interted []string) {
	dao.Data.Store(fmt.Sprint(invertedKey), interted)
}

func (dao *invertedResumeDao) Get(id string) ([]string, bool) {
	if info, ok := dao.Data.Load(id); !ok {
		return nil, ok
	} else {
		return info.([]string), ok
	}
}

func (dao *invertedResumeDao) Del(id string) {
	dao.Data.Delete(id)
}

func (dao *invertedResumeDao) Reload(key string) {

	maps := redis.GetRedis().HGetAll(context.Background(), key).Val()
	for MapKey, value := range maps {
		interted := strings.Split(value, ",")
		if len(interted) == 0 {
			dao.Del(MapKey)
		} else {
			dao.Set(MapKey, interted)
		}
	}
	log.Info().Msgf("job倒排更新成功. %v", key)
}
