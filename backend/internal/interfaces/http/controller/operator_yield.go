package controller

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func (c *controller) ExportOperatorYield(
	ctx context.Context,
	request openapi.ExportOperatorYieldRequestObject,
) (openapi.ExportOperatorYieldResponseObject, error) {
	// openapi_types.Date is forwarded as-is; the usecase/repository reads only Y/M/D
	// and re-interprets the calendar date in JST.
	filter := usecase.OperatorYieldExportFilter{
		DateFrom:   request.Params.DateFrom.Time,
		DateTo:     request.Params.DateTo.Time,
		LocationID: request.Params.LocationId,
		TaskID:     request.Params.TaskId,
		UserID:     request.Params.UserId,
	}

	csvBytes, err := c.operatorYieldExportUsecase.Export(ctx, filter)
	if err != nil {
		return nil, err
	}

	filename := usecase.OperatorYieldExportFilename(time.Now())

	return openapi.ExportOperatorYield200TextcsvResponse{
		Body: bytes.NewReader(csvBytes),
		Headers: openapi.ExportOperatorYield200ResponseHeaders{
			ContentDisposition: fmt.Sprintf(`attachment; filename="%s"`, filename),
		},
	}, nil
}
