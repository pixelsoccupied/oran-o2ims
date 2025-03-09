package alertmanager

import (
	"log/slog"
	"maps"
	"time"

	"github.com/google/uuid"
	api "github.com/openshift-kni/oran-o2ims/internal/service/alarms/api/generated"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/db/models"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/infrastructure"
)

// ConvertAmToAlarmEventRecordModels get alarmEventRecords based on the alertmanager notification and AlarmDefinition
func ConvertAmToAlarmEventRecordModels(alerts *[]api.Alert, infrastructureClient infrastructure.Client) []models.AlarmEventRecord {
	records := make([]models.AlarmEventRecord, 0, len(*alerts))
	for _, alert := range *alerts {
		record := models.AlarmEventRecord{}

		// Validate startsAt is always there
		if alert.StartsAt != nil && !alert.EndsAt.IsZero() {
			record.AlarmRaisedTime = *alert.StartsAt
		} else {
			slog.Error("Alert StartsAt is required, skipping.", "alert", alert)
			continue
		}

		// Validate status is always there, also derive the PerceivedSeverity as needed
		if alert.Status != nil {
			record.AlarmStatus = string(*alert.Status)
			// Make sure the current payload has the right severity
			if *alert.Status == api.Resolved {
				record.PerceivedSeverity = severityToPerceivedSeverity("cleared")
			} else {
				ps, _ := GetPerceivedSeverity(*alert.Labels)
				record.PerceivedSeverity = ps
			}
		} else {
			slog.Error("Alert Status is required, skipping.", "alert", alert)
			continue
		}

		// Validate Fingerprint is always there
		if alert.Fingerprint != nil {
			record.Fingerprint = *alert.Fingerprint
		} else {
			slog.Error("Alert Fingerprint is required, skipping.", "alert", alert)
			continue
		}

		if alert.EndsAt != nil && !alert.EndsAt.IsZero() {
			record.AlarmClearedTime = alert.EndsAt
		}

		// Update Extensions with things we didn't really process
		record.Extensions = getExtensions(*alert.Labels, *alert.Annotations)

		// for caas alerts object is the cluster ID
		record.ObjectID = GetClusterID(*alert.Labels)

		// derive ObjectTypeID from ObjectID
		if record.ObjectID != nil {
			objectTypeID, err := infrastructureClient.GetObjectTypeID(*record.ObjectID)
			if err != nil {
				slog.Warn("Could not get object type ID", "objectID", record.ObjectID, "err", err.Error())
			} else {
				record.ObjectTypeID = &objectTypeID
			}
		}

		// See if possible to pick up additional info from its definition
		if record.ObjectTypeID != nil {
			_, severity := GetPerceivedSeverity(*alert.Labels)
			alarmDefinitionID, err := infrastructureClient.GetAlarmDefinitionID(*record.ObjectTypeID, GetAlertName(*alert.Labels), severity)
			if err != nil {
				slog.Warn("Could not get alarm definition ID", "objectTypeID", *record.ObjectTypeID, "name", GetAlertName(*alert.Labels), "severity", severity, "err", err.Error())
			} else {
				record.AlarmDefinitionID = &alarmDefinitionID
			}
		}

		// Anything else that's not mentioned explicitly will be handled by DB such ID generation and default values as needed.
		records = append(records, record)
	}

	return records
}

func GetClusterID(labels map[string]string) *uuid.UUID {
	val, ok := labels["managed_cluster"]
	if !ok {
		slog.Warn("Could not find managed_cluster", "labels", labels)
		return nil
	}

	id, err := uuid.Parse(val)
	if err != nil {
		slog.Warn("Could convert managed_cluster string to uuid", "labels", labels, "err", err.Error())
		return nil
	}

	return &id
}

// GetAlertName extract name from alert label
func GetAlertName(labels map[string]string) string {
	val, ok := labels["alertname"]
	if !ok {
		// this may never execute but keeping a check just in case
		slog.Warn("Could not find alertname", "labels", labels)
		return "Unknown"
	}

	return val
}

// GetPerceivedSeverity am's `severity` to oran's PerceivedSeverity
func GetPerceivedSeverity(labels map[string]string) (api.PerceivedSeverity, string) {
	severity, ok := labels["severity"]
	if !ok {
		slog.Warn("Could not find severity label", "labels", labels)
		return api.INDETERMINATE, ""
	}

	return severityToPerceivedSeverity(severity), severity
}

func severityToPerceivedSeverity(input string) api.PerceivedSeverity {
	switch input {
	case "cleared":
		return api.CLEARED
	case "critical":
		return api.CRITICAL
	case "major":
		return api.MAJOR
	case "minor", "low":
		return api.MINOR
	case "warning", "info":
		return api.WARNING
	default:
		return api.INDETERMINATE
	}
}

// getExtensions extract oran extension from alert. For caas it's basically the labels and annotations from payload.
func getExtensions(labels, annotations map[string]string) map[string]string {
	if labels == nil {
		labels = make(map[string]string)
	}
	if annotations == nil {
		annotations = make(map[string]string)
	}

	result := make(map[string]string)
	maps.Copy(result, labels)
	maps.Copy(result, annotations)
	return result
}

// ConvertAPIAlertsToWebhook converts a slice of APIAlert objects into a slice of WebhookAlert.
func ConvertAPIAlertsToWebhook(apiAlerts []APIAlert) ([]api.Alert, error) {
	// Handle empty input array
	if len(apiAlerts) == 0 {
		return []api.Alert{}, nil
	}

	webhookAlerts := make([]api.Alert, 0, len(apiAlerts))
	now := time.Now().UTC()

	for _, a := range apiAlerts {
		// Determine the alert status based on endsAt compared to current time.
		//
		// This is strange but API call will always have an `endsAt` regardless if an alert actually "ended" or resolved.
		// AM api `endsAt` is either "now() + resolve_timeout" (future) or time in the past indicating an alert is resolved.
		//
		// A resolved alerts (i.e endsAt from past) is set to be removed from the /alerts but depending on when we call the API
		// we may still see "resolved" alerts, which is why this need to be handled correctly.
		var state api.AlertmanagerNotificationStatus
		var finalEndsAt *time.Time
		if now.Before(a.EndsAt) {
			state = api.Firing
			finalEndsAt = nil
		} else {
			state = api.Resolved
			finalEndsAt = &a.EndsAt
		}

		// Other than endsAt and status (which needs to be computed) everything else is a pass-through
		// to make API payload to a webhook payload
		webhookAlert := api.Alert{
			Annotations:  &a.Annotations,
			Labels:       &a.Labels,
			StartsAt:     &a.StartsAt,
			EndsAt:       finalEndsAt,
			Fingerprint:  &a.Fingerprint,
			GeneratorURL: &a.GeneratorURL,
			Status:       &state,
		}

		webhookAlerts = append(webhookAlerts, webhookAlert)
	}

	slog.Info("Converted from API to Webhook", "alerts", len(webhookAlerts))
	return webhookAlerts, nil
}
