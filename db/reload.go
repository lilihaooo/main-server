/*
   @Author: ken.quan<wallestore@hotmail.com>
   @Desc: 说明
   @Create At: 2021/4/14
*/
package db

import (
	"context"
	"github.com/rs/zerolog/log"
	"gitlab.com/canyinxinxi/main-server/db/inverted_job"
	"gitlab.com/canyinxinxi/main-server/db/inverted_resume"
	"gitlab.com/canyinxinxi/main-server/db/job"
	"gitlab.com/canyinxinxi/main-server/db/resume"
	"gitlab.com/canyinxinxi/main-server/db/user"
	"gitlab.com/canyinxinxi/main-server/module/redis"
	"gitlab.com/canyinxinxi/main-server/pb/index"
	"strings"
)

const (
	USER_LIST       = "a:index:forward.user.list"
	Job_LIST        = "a:index:forward.job.list"
	RESUME_LIST     = "a:index:forward.resume.list"
	INVERTED_JOB    = "a:index:inverted.job"
	INVERTED_RESUME = "a:index:inverted.resume"
)

func Reload() (err error) {
	ctx := context.Background()
	//加载用户
	userKeys := redis.GetRedis().SMembers(ctx, USER_LIST).Val()
	user.GetUserDao().Reload(userKeys)
	//加载job
	jobKeys := redis.GetRedis().SMembers(ctx, Job_LIST).Val()
	job.GetJobDao().Reload(jobKeys)

	//加载resume
	resumeKeys := redis.GetRedis().SMembers(ctx, RESUME_LIST).Val()
	resume.GetResumeDao().Reload(resumeKeys)

	//job 倒排加载
	inverted_job.GetInvertedJobDao().Reload(INVERTED_JOB)
	//resume 倒排加载
	inverted_resume.GetInvertedResumeDao().Reload(INVERTED_RESUME)

	log.Info().Msg("全量数据同步完成")
	return
}

func Update(notify *index.UpdateNotify) {
	ctx := context.Background()
	switch notify.Transaction.Type {
	case index.UpdateNotify_User:
		user.GetUserDao().Reload(notify.Transaction.Keys)
	case index.UpdateNotify_Job:
		job.GetJobDao().Reload(notify.Transaction.Keys)
	case index.UpdateNotify_Resume:
		resume.GetResumeDao().Reload(notify.Transaction.Keys)

	case index.UpdateNotify_INVERTED_Job:
		for _, id := range notify.Transaction.Keys {
			v := redis.GetRedis().HMGet(ctx, INVERTED_JOB, id).String()
			interted := strings.Split(v, ",")
			inverted_job.GetInvertedJobDao().Set(id, interted)
		}
		log.Info().Msgf("倒排索引Job更新成功. %v", notify.Transaction.Keys)
	case index.UpdateNotify_INVERTED_Resume:
		for _, id := range notify.Transaction.Keys {
			v := redis.GetRedis().HMGet(ctx, INVERTED_RESUME, id).String()
			interted := strings.Split(v, ",")
			inverted_resume.GetInvertedResumeDao().Set(id, interted)
		}
		log.Info().Msgf("倒排索Resume更新成功. %v", notify.Transaction.Keys)

	}
}
