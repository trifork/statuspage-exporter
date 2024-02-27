package engines

import (
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sergeyshevch/statuspage-exporter/pkg/config"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/statusio"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/statuspageio"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/types"
)

var errUnknownStatusPageType = errors.New("unknown statuspage type")

// FetchStatus detect statuspage type and fetch its status.
func FetchStatus(
	config *config.ExporterConfig,
	targetURL string,
	componentStatus *prometheus.GaugeVec,
	overallStatus *prometheus.GaugeVec,
) error {
	restyClient := resty.New().
		EnableTrace().
		SetTimeout(config.ClientTimeout).
		SetRetryCount(config.RetryCount)

	detector := NewDetector()

	statusPageType := detector.DetectStatusPageType(config.Log, restyClient, targetURL)
	if statusPageType == types.UnknownType {
		return errUnknownStatusPageType
	}

	switch statusPageType {
	case types.StatusPageIOType:
		return statuspageio.FetchStatusPage(
			config.Log,
			targetURL,
			restyClient,
			componentStatus,
			overallStatus,
		)
	case types.StatusIOType:
		return statusio.FetchStatusPage(config.Log, targetURL, restyClient, componentStatus, overallStatus)
	case types.UnknownType:
		return errUnknownStatusPageType
	default:
		return errUnknownStatusPageType
	}
}
