package alertmanager

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"

	api "github.com/openshift-kni/oran-o2ims/internal/service/alarms/api/generated"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/db/repo"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/infrastructure"
)

// Source of Alertmanager payload either through API or Webhook
type Source string

const (
	API     Source = "API"
	Webhook Source = "Webhook"
)

// HandleAlerts can be called when a payload from Webhook or API `/alerts` is received
// Webhook is our primary and API as our backup and sync mechanism
func HandleAlerts(ctx context.Context, clients []infrastructure.Client, repository repo.AlarmRepositoryInterface, alerts *[]api.Alert, source Source) error {
	if len(*alerts) == 0 {
		return nil
	}

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
	if err := repository.WithTransaction(ctx, func(tx pgx.Tx) error { //nolint:wrapcheck
		slog.Info("Handling alerts", "source", string(source))
		// genID to determine if stale
		generationID := time.Now().UnixNano()

		// Insert or update with alerts
		if err := repository.UpsertAlarmEventRecord(ctx, aerModels, generationID); err != nil {
			return fmt.Errorf("failed to upsert alarm event record model: %w", err)
		}

		// Resolve stale only if source is API since `/alerts` as this step only works if we have full set of alerts
		if source == API {
			if err := repository.ResolveStaleAlarmEventRecord(ctx, int(generationID)); err != nil {
				return fmt.Errorf("could not resolve notification: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to handle alerts from %s: %w", string(source), err)
	}

	slog.Info("Successfully handled AlarmEventRecords", "source", string(source))
	return nil
}
