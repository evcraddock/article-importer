package tasks

import (
	"github.com/evcraddock/article-importer/config"
	"github.com/evcraddock/article-importer/service"
)

type Task struct {
	service 			*service.HttpService
}

func NewTask(settings *config.Settings) *Task {
	service := service.NewHttpService(settings)

	task := &Task{
		service,
	}

	return task
}