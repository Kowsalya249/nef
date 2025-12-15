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

// PostAfSessionQos creates an AF Session with QoS and relays it to PCF.
func (p *Processor) PostAfSessionQos(c *gin.Context, afID string, asc *models.AppSessionContext) {
	nefCtx := p.Context()
	af := nefCtx.GetAf(afID)
	if af == nil {
		af = nefCtx.NewAf(afID)
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

	sessID := uuid.New().String()
	qosSess := &context.AfQosSession{
		SessID:      sessID,
		AppSessID:   appSessID,
		NotifCorrID: corrID,
		Payload:     asc,
		Log:         af.Log.WithField(logger.FieldSubID, sessID),
	}
	af.AddQosSession(qosSess)
	nefCtx.AddAf(af)

	self := p.genAfSessionQosURI(afID, sessID)
	headers := map[string][]string{
		"Location": {self},
	}
	for hdrName, hdrValues := range headers {
		for _, hdrValue := range hdrValues {
			c.Header(hdrName, hdrValue)
		}
	}
	asc.Self = self
	c.JSON(http.StatusCreated, asc)
}

// GetAfSessionQos returns a stored AF QoS session representation.
func (p *Processor) GetAfSessionQos(c *gin.Context, afID, sessID string) {
	af := p.Context().GetAf(afID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("AF is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}
	af.Mu.RLock()
	defer af.Mu.RUnlock()

	qosSess, ok := af.GetQosSession(sessID)
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("QoS session is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}
	c.JSON(http.StatusOK, qosSess.Payload)
}

// PatchAfSessionQos updates an AF QoS session and relays to PCF.
func (p *Processor) PatchAfSessionQos(c *gin.Context, afID, sessID string, ascUpdate *models.AppSessionContextUpdateData) {
	af := p.Context().GetAf(afID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("AF is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	qosSess, ok := af.GetQosSession(sessID)
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("QoS session is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	respAsc, pd, err := p.Consumer().PatchAppSession(qosSess.AppSessID, ascUpdate)
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
		qosSess.Payload = respAsc
		c.JSON(http.StatusOK, respAsc)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteAfSessionQos deletes an AF QoS session and relays to PCF.
func (p *Processor) DeleteAfSessionQos(c *gin.Context, afID, sessID string) {
	af := p.Context().GetAf(afID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("AF is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	qosSess, ok := af.GetQosSession(sessID)
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("QoS session is not found")
		c.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		c.JSON(int(pd.Status), pd)
		return
	}

	rspCode, pd, err := p.Consumer().DeleteAppSession(qosSess.AppSessID)
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

	af.DeleteQosSession(sessID)

	if rspCode == http.StatusNoContent || rspCode == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.Status(rspCode)
}

func (p *Processor) genAfSessionQosURI(afID, sessID string) string {
	return factory.AfSessionQosResUriPrefix + "/" + afID + "/sessions/" + sessID
}
