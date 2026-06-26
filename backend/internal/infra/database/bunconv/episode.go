package bunconv

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
)

func EntityToEpisodeModel(ee entity.Episode, et entity.TaskVersion) model.Episode {
	ep := model.Episode{
		ID:             ee.ID,
		IDNatural:      ee.IDNatural,
		OrganizationID: ee.OrganizationID,
		TaskID:         et.TaskID,
		TaskVersionID:  ee.TaskVersionID,
		LocationID:     ee.LocationID,
		RobotID:        ee.RobotID,
		UserID:         ee.UserID,
		RecordedByID:   ee.RecordedByID,
		StartedAt:      ee.StartedAt,
		FinishedAt:     ee.FinishedAt,
		Status:         ee.CollectionStatus,
		ErrorDetails:   ee.ErrorDetails,
		CreatedAt:      ee.CreatedAt,
		UpdatedAt:      &ee.UpdatedAt,
	}
	if len(ee.ParameterValues) > 0 {
		var pv map[string]string
		if err := json.Unmarshal(ee.ParameterValues, &pv); err != nil {
			log.Error().Err(err).Str("id_natural", ee.IDNatural).Str("raw", string(ee.ParameterValues)).Msg("failed to unmarshal episode parameter values")
		} else {
			ep.ParameterValues = pv
		}
	}
	return ep
}

func ParameterValuesToJSON(pv map[string]string) json.RawMessage {
	if len(pv) == 0 {
		return nil
	}
	b, err := json.Marshal(pv)
	if err != nil {
		return nil
	}
	return b
}
