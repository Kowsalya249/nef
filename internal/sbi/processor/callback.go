package processor

import (
	"net/http"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/metrics/sbi"
	"github.com/gin-gonic/gin"
)

func (p *Processor) SmfNotification(
	c *gin.Context,
	eeNotif *models.NsmfEventExposureNotification,
) {
	logger.TrafInfluLog.Infof("SmfNotification - NotifId[%s]", eeNotif.NotifId)

	af, sub := p.Context().FindAfSub(eeNotif.NotifId)
	if sub == nil {
		pd := openapi.ProblemDetailsDataNotFound("Subscription is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(http.StatusNotFound, pd)
		return
	}

	af.Mu.RLock()
	defer af.Mu.RUnlock()

	// TODO: Notify AF

	c.JSON(http.StatusOK, nil)
}

func (p *Processor) AfSessionQosNotification(
	c *gin.Context,
	corrID string,
	ascUpdate *models.AppSessionContextUpdateData,
) {
	logger.TrafInfluLog.Infof("AfSessionQosNotification - CorrID[%s]", corrID)

	af, sess := p.Context().FindAfQosSessionByCorrID(corrID)
	if sess == nil {
		pd := openapi.ProblemDetailsDataNotFound("QoS session is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(http.StatusNotFound, pd)
		return
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	// TODO: Forward notification to AF endpoint if mediation is required.
	// For now, just log and update stored payload if PCF returned a full context.
	if ascUpdate != nil && sess.Payload != nil && ascUpdate.AscReqData != nil {
		// Minimal merge: replace AscReqData with latest from PCF notification.
		sess.Payload.AscReqData = ascUpdate.AscReqData
	}

	c.Status(http.StatusNoContent)
}
