package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/free5gc/nef/internal/context"
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
<<<<<<< Updated upstream
=======

func (p *Processor) AsSessionQosNotification(
	c *gin.Context,
	corrID string,
	ascUpdate *models.AppSessionContextUpdateData,
) {
	logger.TrafInfluLog.Infof("AsSessionQosNotification - CorrID[%s]", corrID)

	af, sub := p.Context().FindAfQosSubscriptionByCorrID(corrID)
	if sub == nil {
		pd := openapi.ProblemDetailsDataNotFound("QoS subscription is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(http.StatusNotFound, pd)
		return
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	if ascUpdate != nil {
		sub.LastUpdate = ascUpdate
	}

	if err := p.forwardAsSessionQosNotification(sub, corrID, ascUpdate); err != nil {
		logger.TrafInfluLog.Warnf("Failed to forward QoS notification to AF: %v", err)
	}

	c.Status(http.StatusNoContent)
}

func (p *Processor) forwardAsSessionQosNotification(
	sub *context.AfQosSubscription,
	corrID string,
	ascUpdate *models.AppSessionContextUpdateData,
) error {
	if sub == nil || sub.Payload == nil || sub.Payload.AscReqData == nil {
		return fmt.Errorf("subscription or request data missing for forwarding")
	}

	dest := sub.Payload.AscReqData.NotifUri
	if dest == "" {
		return fmt.Errorf("notification URI missing in ascReqData")
	}

	body, err := json.Marshal(ascUpdate)
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, dest, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Correlation-Id", corrID)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("call AF notification endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("AF notification endpoint returned %d", resp.StatusCode)
	}
	return nil
}
>>>>>>> Stashed changes
