package controller

import (
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/dto"
	"github.com/HEBNUOJ/model"
	"github.com/HEBNUOJ/response"
	"github.com/HEBNUOJ/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type ProblemController struct{}

func (p *ProblemController) AddProblem(ctx *gin.Context) {
	requestProblem := dto.PublicProblemDto{}
	ctx.Bind(&requestProblem)
	// 获取参数
	title := requestProblem.Title
	description := requestProblem.Description
	input := requestProblem.Input
	output := requestProblem.Output
	sampleInput := requestProblem.SampleInput
	sampleOutput := requestProblem.SampleOutput
	spj := requestProblem.Spj
	hint := requestProblem.Hint
	source := requestProblem.Source
	timeLimit := requestProblem.TimeLimit
	memoryLimit := requestProblem.MemoryLimit
	defunct := requestProblem.Defunct
	degree := requestProblem.Degree

	// 参数校检
	errString := ""
	switch {
	case len(title) == 0:
		errString = "标题不能为空"
	case timeLimit < 1000:
		errString = "时间限制至少为1000ms"
	case memoryLimit < 32:
		errString = "内存限制至少为32MB"
	}
	if len(errString) > 0 {
		response.Response(ctx, http.StatusOK, 422, nil, errString)
		return
	}

	newProblem := model.PublicProblem{
		Title:        title,
		Description:  description,
		Input:        input,
		Output:       output,
		SampleInput:  sampleInput,
		SampleOutput: sampleOutput,
		Spj:          spj,
		Hint:         hint,
		Source:       source,
		InDate:       time.Now(),
		TimeLimit:    timeLimit,
		MemoryLimit:  memoryLimit,
		Defunct:      defunct,
		Accepted:     0,
		Submit:       0,
		Degree:       degree,
	}

	common.GetDB().Save(&newProblem)
	response.Success(ctx, nil, "添加题目成功")
}

func (p *ProblemController) UpdateProblem(ctx *gin.Context) {
	p.AddProblem(ctx)
}

func (p *ProblemController) DelProblem(ctx *gin.Context) {
	requestProblem := dto.PublicProblemDto{}
	ctx.Bind(&requestProblem)
	var problem model.PublicProblem
	common.GetDB().Where("id = ?", requestProblem.Id).First(&problem)
	if problem.Id == 0 {
		response.Response(ctx, http.StatusOK, 422, nil, "题目不存在")
		return
	}
	common.GetDB().Delete(problem)
	response.Success(ctx, nil, "删除成功")
}

func (p *ProblemController) QueryProblem(ctx *gin.Context) {
	id := ctx.Param("id")
	var problem model.PublicProblem
	common.GetDB().Where("id = ?", id).First(&problem)
	if problem.Id == 0 {
		response.Response(ctx, http.StatusOK, 422, nil, "题目不存在")
		return
	}
	response.Success(ctx, gin.H{"problem": dto.ToProblemDto(problem)}, "查询成功")
}

func (p *ProblemController) ShowProblemList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	const pageSize = 40
	db := common.GetDB()
	if page > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	problems := make([]model.PublicProblem, 0)
	if err := db.Find(&problems).Error; err != nil {
		utils.Log("problemset.log", 1).Println("分页查询失败", err)
		return
	}
	problemJson := make([]dto.PublicProblemDto, 0)
	for i := 0; i < len(problems); i++ {
		problemJson[i] = dto.ToProblemDto(problems[i])
	}

	response.Success(ctx, gin.H{"problems": problemJson}, "分页查询成功")
}

func (p *ProblemController) SubmitProblem(ctx *gin.Context) {

}
