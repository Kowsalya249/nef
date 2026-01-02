package processor

import (
	"net/http"

	"github.com/free5gc/nef/internal/context"
	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/nef/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/metrics/sbi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListAsSessionQosSubs returns all QoS subscriptions for an SCS/AS.
func (p *Processor) ListAsSessionQosSubs(c *gin.Context, scsAsID string) {
	af := p.Context().GetAf(scsAsID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("SCS/AS is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	af.Mu.RLock()
	defer af.Mu.RUnlock()

	var subs []models.AppSessionContext
	for _, sub := range af.QosSubs {
		if sub.Payload != nil {
			subs = append(subs, *sub.Payload)
		}
	}
	c.JSON(http.StatusOK, subs)
}

// PostAsSessionQosSub creates a QoS subscription and relays it to PCF.
func (p *Processor) PostAsSessionQosSub(c *gin.Context, scsAsID string, asc *models.AppSessionContext) {
	nefCtx := p.Context()
	af := nefCtx.GetAf(scsAsID)
	if af == nil {
		af = nefCtx.NewAf(scsAsID)
		if af == nil {
			pd := openapi.ProblemDetailsSystemFailure("No resource can be allocated")
			c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
			c.JSON(int(pd.Status), pd)
			return
		}
	}

	corrID := uuid.New().String()

	appSessID, pd, err := p.Consumer().PostAppSessions(asc)
	switch {
	case pd != nil:
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	case err != nil:
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, problemDetails.Cause)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	subID := uuid.New().String()
	qosSub := &context.AfQosSubscription{
		SubscriptionID: subID,
		AppSessID:      appSessID,
		NotifCorrID:    corrID,
		Payload:        asc,
		Log:            af.Log.WithField(logger.FieldSubID, subID),
	}
	af.AddQosSubscription(qosSub)
	nefCtx.AddAf(af)

	self := p.genAsSessionQosURI(scsAsID, subID)
	headers := map[string][]string{
		"Location": {self},
	}
	for hdrName, hdrValues := range headers {
		for _, hdrValue := range hdrValues {
			c.Header(hdrName, hdrValue)
		}
	}
	c.JSON(http.StatusCreated, asc)
}

// GetAsSessionQosSub returns a stored QoS subscription representation.
func (p *Processor) GetAsSessionQosSub(c *gin.Context, scsAsID, subID string) {
	af := p.Context().GetAf(scsAsID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("SCS/AS is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}
	af.Mu.RLock()
	defer af.Mu.RUnlock()

	qosSub, ok := af.GetQosSubscription(subID)
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("QoS subscription is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}
	c.JSON(http.StatusOK, qosSub.Payload)
}

// PutAsSessionQosSub updates a QoS subscription (idempotent) and relays to PCF.
func (p *Processor) PutAsSessionQosSub(
	c *gin.Context,
	scsAsID, subID string,
	ascUpdate *models.AppSessionContextUpdateData,
) {
	p.updateAsSessionQosSub(c, scsAsID, subID, ascUpdate)
}

// PatchAsSessionQosSub updates a QoS subscription (partial) and relays to PCF.
func (p *Processor) PatchAsSessionQosSub(
	c *gin.Context,
	scsAsID, subID string,
	ascUpdate *models.AppSessionContextUpdateData,
) {
	p.updateAsSessionQosSub(c, scsAsID, subID, ascUpdate)
}

func (p *Processor) updateAsSessionQosSub(
	c *gin.Context,
	scsAsID, subID string,
	ascUpdate *models.AppSessionContextUpdateData,
) {
	af := p.Context().GetAf(scsAsID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("SCS/AS is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	qosSub, ok := af.GetQosSubscription(subID)
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("QoS subscription is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	respAsc, pd, err := p.Consumer().PatchAppSession(qosSub.AppSessID, ascUpdate)
	switch {
	case pd != nil:
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	case err != nil:
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, problemDetails.Cause)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	if respAsc != nil {
		qosSub.Payload = respAsc
		c.JSON(http.StatusOK, respAsc)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteAsSessionQosSub deletes a QoS subscription and relays to PCF.
func (p *Processor) DeleteAsSessionQosSub(c *gin.Context, scsAsID, subID string) {
	af := p.Context().GetAf(scsAsID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("SCS/AS is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	qosSub, ok := af.GetQosSubscription(subID)
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("QoS subscription is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	rspCode, pd, err := p.Consumer().DeleteAppSession(qosSub.AppSessID)
	switch {
	case pd != nil:
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	case err != nil:
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, problemDetails.Cause)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	af.DeleteQosSubscription(subID)

	if rspCode == http.StatusNoContent || rspCode == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.Status(rspCode)
}

func (p *Processor) genAsSessionQosURI(scsAsID, subID string) string {
	return factory.AsSessionQosResUriPrefix + "/" + scsAsID + "/subscriptions/" + subID
}
