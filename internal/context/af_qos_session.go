package context

import (
	"github.com/free5gc/openapi/models"
	"github.com/sirupsen/logrus"
)

// AfQosSubscription represents a QoS exposure subscription tracked by NEF.
type AfQosSubscription struct {
	SubscriptionID string
	AppSessID      string
	NotifCorrID    string
	Payload        *models.AppSessionContext
	LastUpdate     *models.AppSessionContextUpdateData
	Log            *logrus.Entry
}
