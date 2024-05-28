package utils

import "github.com/fernandonogueira/statuspage-exporter/pkg/engines/types"

func StatusToString(status types.Status) string {
	switch status {
	case types.UnknownStatus:
		return "unknown"
	case types.OperationalStatus:
		return "operational"
	case types.PlannedMaintenanceStatus:
		return "planned_maintenance"
	case types.DegradedPerformanceStatus:
		return "degraded_performance"
	case types.PartialOutageStatus:
		return "partial_outage"
	case types.MajorOutageStatus:
		return "major_outage"
	case types.SecurityIssueStatus:
		return "security_issue"
	default:
		return "unknown"
	}
}
