package repo

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/db/models"
	commonmodels "github.com/openshift-kni/oran-o2ims/internal/service/common/db/models"
)

//go:generate mockgen -source=alarms_repository_interface.go -destination=generated/mock_repo.generated.go -package=generated

type AlarmRepositoryInterface interface {
	GetAlarmEventRecords(ctx context.Context) ([]models.AlarmEventRecord, error)
	PatchAlarmEventRecordACK(ctx context.Context, id uuid.UUID, record *models.AlarmEventRecord) (*models.AlarmEventRecord, error)
	GetAlarmEventRecord(ctx context.Context, id uuid.UUID) (*models.AlarmEventRecord, error)
	CreateServiceConfiguration(ctx context.Context, defaultRetentionPeriod int) (*models.ServiceConfiguration, error)
	GetServiceConfigurations(ctx context.Context) ([]models.ServiceConfiguration, error)
	UpdateServiceConfiguration(ctx context.Context, id uuid.UUID, record *models.ServiceConfiguration) (*models.ServiceConfiguration, error)
	GetAlarmSubscriptions(ctx context.Context) ([]models.AlarmSubscription, error)
	DeleteAlarmSubscription(ctx context.Context, id uuid.UUID) (int64, error)
	CreateAlarmSubscription(ctx context.Context, record models.AlarmSubscription) (*models.AlarmSubscription, error)
	GetAlarmSubscription(ctx context.Context, id uuid.UUID) (*models.AlarmSubscription, error)
	UpsertAlarmEventRecord(ctx context.Context, records []models.AlarmEventRecord, generationID int64) error
	UpdateSubscriptionEventCursor(ctx context.Context, subscription models.AlarmSubscription) error
	GetAllAlarmsDataChange(ctx context.Context) ([]commonmodels.DataChangeEvent, error)
	DeleteAlarmsDataChange(ctx context.Context, dataChangeId uuid.UUID) error
	WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error
	ResolveNotificationWithStaleGenID(ctx context.Context, generationID int) error
}
