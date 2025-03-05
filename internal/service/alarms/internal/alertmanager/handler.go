package alertmanager

import (
	"context"
	"fmt"
	"log/slog"

	api "github.com/openshift-kni/oran-o2ims/internal/service/alarms/api/generated"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/db/repo"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/infrastructure"
)

// HandleAMAlerts this is called from both Webhook and API alerts payload
func HandleAMAlerts(ctx context.Context, clients []infrastructure.Client, repository repo.AlarmRepositoryInterface, alerts *[]api.Alert, generationID int64, fullSync bool) error {
	// Get cached cluster server data
	var clusterServer infrastructure.Client
	for i := range clients {
		if clients[i].Name() == infrastructure.Name {
			clusterServer = clients[i]
		}
	}

	// Combine possible definitions with events
	aerModels := ConvertAmToAlarmEventRecordModels(alerts, clusterServer)

	// Insert and update AlarmEventRecord
	if err := repository.UpsertAlarmEventRecord(ctx, aerModels, generationID, fullSync); err != nil {
		msg := "failed to upsert AlarmEventRecord to db"
		slog.Error(msg, "error", err)
		return fmt.Errorf("%s: %w", msg, err)
	}

	slog.Info("Successfully upserted AlarmEventRecords to db")
	return nil
}
