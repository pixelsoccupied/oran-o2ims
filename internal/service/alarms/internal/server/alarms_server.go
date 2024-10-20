package server

import (
	"context"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/api/generated"
)

type AlarmsServer struct {
}

var _ generated.StrictServerInterface = (*AlarmsServer)(nil)

func (a AlarmsServer) GetSubscriptions(ctx context.Context, request generated.GetSubscriptionsRequestObject) (generated.GetSubscriptionsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) CreateSubscription(ctx context.Context, request generated.CreateSubscriptionRequestObject) (generated.CreateSubscriptionResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) DeleteSubscription(ctx context.Context, request generated.DeleteSubscriptionRequestObject) (generated.DeleteSubscriptionResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) GetSubscription(ctx context.Context, request generated.GetSubscriptionRequestObject) (generated.GetSubscriptionResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) GetAlarms(ctx context.Context, request generated.GetAlarmsRequestObject) (generated.GetAlarmsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) GetAlarm(ctx context.Context, request generated.GetAlarmRequestObject) (generated.GetAlarmResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) AckAlarm(ctx context.Context, request generated.AckAlarmRequestObject) (generated.AckAlarmResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) GetProbableCauses(ctx context.Context, request generated.GetProbableCausesRequestObject) (generated.GetProbableCausesResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) GetProbableCause(ctx context.Context, request generated.GetProbableCauseRequestObject) (generated.GetProbableCauseResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) AmNotification(ctx context.Context, request generated.AmNotificationRequestObject) (generated.AmNotificationResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (a AlarmsServer) HwNotification(ctx context.Context, request generated.HwNotificationRequestObject) (generated.HwNotificationResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
