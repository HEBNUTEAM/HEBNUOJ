package dto

import (
	"github.com/HEBNUOJ/model"
	"time"
)

type PublicProblemDto struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Input        string    `json:"input"`
	Output       string    `json:"output"`
	SampleInput  string    `json:"sampleinput"`
	SampleOutput string    `json:"sampleoutput"`
	Spj          bool      `json:"spj"`
	Hint         string    `json:"hint"`
	Source       string    `json:"source"`
	InDate       time.Time `json:"indate"`
	TimeLimit    int       `json:"timelimit"`
	MemoryLimit  int       `json:"memorylimit"`
	Defunct      int       `json:"defunct"`
	Accepted     int       `json:"accepted"`
	Submit       int       `json:"submit"`
	Degree       string    `json:"degree"`
}

func ToProblemDto(problem model.PublicProblem) PublicProblemDto {
	return PublicProblemDto{
		Id:           problem.Id,
		Title:        problem.Title,
		Description:  problem.Description,
		Input:        problem.Input,
		Output:       problem.Output,
		SampleInput:  problem.SampleInput,
		SampleOutput: problem.SampleOutput,
		Spj:          problem.Spj,
		Hint:         problem.Hint,
		Source:       problem.Source,
		InDate:       problem.InDate,
		TimeLimit:    problem.TimeLimit,
		MemoryLimit:  problem.MemoryLimit,
		Defunct:      problem.Defunct,
		Accepted:     problem.Accepted,
		Submit:       problem.Submit,
		Degree:       problem.Degree,
	}
}
