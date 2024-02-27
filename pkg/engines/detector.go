package engines

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/statusio"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/statuspageio"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/types"
)

type Detector struct {
	StatusPageTypesBuffer map[string]types.EngineType
}

func NewDetector() *Detector {
	return &Detector{
		StatusPageTypesBuffer: make(map[string]types.EngineType),
	}
}

// DetectStatusPageType detects statuspage engine for given statuspage URLs.
func (d *Detector) DetectStatusPageType(
	log *zap.Logger,
	restyClient *resty.Client,
	targetURL string,
) types.EngineType {
	if engine, ok := d.StatusPageTypesBuffer[targetURL]; ok {
		return engine
	}

	if statuspageio.IsStatusPageIOPage(log, targetURL, restyClient) {
		log.Info("Detected StatusPage.io page", zap.String("url", targetURL))

		d.StatusPageTypesBuffer[targetURL] = types.StatusPageIOType

		return types.StatusPageIOType
	} else if statusio.IsStatusIOPage(log, targetURL, restyClient) {
		log.Info("Detected Status.io page", zap.String("url", targetURL))

		d.StatusPageTypesBuffer[targetURL] = types.StatusIOType

		return types.StatusIOType
	}

	return types.UnknownType
}
