package bunconv

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
)

func EntityToTaskVersionModel(e entity.TaskVersion) *model.TaskVersion {
	tv := &model.TaskVersion{
		ID:                              e.ID,
		IDNatural:                       e.IDNatural,
		OrganizationID:                  e.OrganizationID,
		TaskID:                          e.TaskID,
		Version:                         e.Version,
		DisplayName:                     e.DisplayName,
		ApprovalStatus:                  e.ApprovalStatus,
		CreatedAt:                       e.CreatedAt,
		TargetDurationSeconds:           e.TargetDurationSeconds,
		TargetEpisodeCount:              e.TargetEpisodeCount,
		TargetDurationPerEpisodeSeconds: e.TargetDurationPerEpisodeSeconds,
	}
	if len(e.Parameters) > 0 {
		var params []model.TaskVersionParameter
		if err := json.Unmarshal(e.Parameters, &params); err != nil {
			log.Error().Err(err).Str("id_natural", e.IDNatural).Str("raw", string(e.Parameters)).Msg("failed to unmarshal task version parameters")
		} else {
			tv.Parameters = params
		}
	}
	return tv
}

func TaskVersionParametersToJSON(params []model.TaskVersionParameter) json.RawMessage {
	if len(params) == 0 {
		return nil
	}
	b, err := json.Marshal(params)
	if err != nil {
		return nil
	}
	return b
}

func EntitiesToTaskVersionModels(entities []entity.TaskVersion) model.TaskVersions {
	result := make(model.TaskVersions, 0, len(entities))
	for _, e := range entities {
		result = append(result, EntityToTaskVersionModel(e))
	}
	return result
}
