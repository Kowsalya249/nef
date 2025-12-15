package context

import (
	"github.com/free5gc/openapi/models"
	"github.com/sirupsen/logrus"
)

// AfQosSession represents an AF-originated QoS session tracked by NEF.
type AfQosSession struct {
	SessID      string
	AppSessID   string
	NotifCorrID string
	Payload     *models.AppSessionContext
	Log         *logrus.Entry
}
