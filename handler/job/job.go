package job

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/canyinxinxi/main-server/db/job"
	"gitlab.com/canyinxinxi/main-server/global"
	"gitlab.com/canyinxinxi/main-server/handler/switch_context"
	"gitlab.com/canyinxinxi/main-server/module/producet"
	"gitlab.com/canyinxinxi/main-server/module/redis"
	"gitlab.com/canyinxinxi/main-server/pb/index"
	job2 "gitlab.com/canyinxinxi/main-server/pb/job"
	"net/http"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/25
*@description:
 */
func Create(ctx *gin.Context) {
	requestInfo := new(index.Job)
	responseInfo := new(job2.JobResponse)
	switch_context.SwitchContextForRequest(ctx, requestInfo)
	CreatJob(requestInfo)
	requestInfo.Status = global.Status_Examine
	job.GetJobDao().Set(requestInfo)
	responseInfo.Code = http.StatusOK
	responseInfo.Job = requestInfo
	//发送job日志
	producet.SendJobRr(requestInfo)
	switch_context.SwitchContextForResponse(ctx, responseInfo)
	return

}
func Info(ctx *gin.Context) {
	id := ctx.Param("jobId")
	responseInfo := new(job2.JobResponse)
	jobInfo, ok := job.GetJobDao().Get(id)
	if ok {
		responseInfo.Code = http.StatusOK
		responseInfo.Job = jobInfo
	} else {
		responseInfo.Code = http.StatusNoContent
		responseInfo.Message = "信息不存在"

	}
	switch_context.SwitchContextForResponse(ctx, responseInfo)
	return

}
func MyJob(ctx *gin.Context) {
	responseInfo := new(job2.JobSResponse)
	jobIds := job.GetJobDao().GetMy(ctx.Request.Header.Get("user_id"))
	var jobs []*index.Job
	for _, id := range jobIds {
		jobObject, ok := job.GetJobDao().Get(id)
		if ok {
			jobs = append(jobs, jobObject)
		}
	}
	responseInfo.Code = http.StatusOK
	responseInfo.Jobs = jobs

	switch_context.SwitchContextForResponse(ctx, responseInfo)
	return

}
func Update(ctx *gin.Context) {
	requestInfo := new(index.Job)
	responseInfo := new(job2.JobResponse)
	switch_context.SwitchContextForRequest(ctx, requestInfo)
	if requestInfo.Id <= 0 {
		CreatJob(requestInfo)
	}
	requestInfo.Status = global.Status_Examine
	job.GetJobDao().Set(requestInfo)
	//发送job日志
	producet.SendJobRr(requestInfo)
	responseInfo.Code = http.StatusOK
	responseInfo.Job = requestInfo
	switch_context.SwitchContextForResponse(ctx, responseInfo)
	return
}

func CreatJob(job *index.Job) {
	job.Id = redis.GetMaxJobId()
}
