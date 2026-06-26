package entity

import (
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
)

func TestEpisodeSubTaskUsesDomainTaskResult(t *testing.T) {
	episodeSubTask := EpisodeSubTask{
		TaskResult: model.TaskResultUndetermined,
	}

	var _ model.TaskResult = episodeSubTask.TaskResult
}
