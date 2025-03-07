package alertmanager

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	api "github.com/openshift-kni/oran-o2ims/internal/service/alarms/api/generated"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/db/repo"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/infrastructure"
)

type Source string

// Define constants for your different sources
const (
	API     Source = "API"
	Webhook Source = "Webhook"
)

// HandleAMAlerts this is called from both Webhook and API alerts payload
func HandleAMAlerts(ctx context.Context, clients []infrastructure.Client, repository repo.AlarmRepositoryInterface, alerts *[]api.Alert, generationID int64, source Source) error {
	// Get cached cluster server data
	var clusterServer infrastructure.Client
	for i := range clients {
		if clients[i].Name() == infrastructure.Name {
			clusterServer = clients[i]
		}
	}

	// Combine possible definitions with events
	aerModels := ConvertAmToAlarmEventRecordModels(alerts, clusterServer)

	// Insert and update AlarmEventRecord and optionally resolve stale
	return repository.WithTransaction(ctx, func(tx pgx.Tx) error { //nolint:wrapcheck
		if err := repository.UpsertAlarmEventRecord(ctx, aerModels, generationID); err != nil {
			return fmt.Errorf("failed to upsert alarm event record model: %w", err)
		}

		if source == API { // Resolve stale only if source is API since `/alerts` endpoint gets us the full set of alerts
			if err := repository.ResolveNotificationWithStaleGenID(ctx, int(generationID)); err != nil {
				return fmt.Errorf("could not resolve notification: %w", err)
			}
		}

		slog.Info("Successfully handled AlarmEventRecords")
		return nil
	})
}
