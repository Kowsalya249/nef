package sbi

import (
	"net/http"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/metrics/sbi"
	"github.com/gin-gonic/gin"
)

func (s *Server) getAfSessionQosRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodPost,
			Pattern: "/:afId/sessions",
			APIFunc: s.apiPostAfSessionQos,
		},
		{
			Method:  http.MethodGet,
			Pattern: "/:afId/sessions/:sessId",
			APIFunc: s.apiGetAfSessionQos,
		},
		{
			Method:  http.MethodPatch,
			Pattern: "/:afId/sessions/:sessId",
			APIFunc: s.apiPatchAfSessionQos,
		},
		{
			Method:  http.MethodDelete,
			Pattern: "/:afId/sessions/:sessId",
			APIFunc: s.apiDeleteAfSessionQos,
		},
	}
}

func (s *Server) apiPostAfSessionQos(gc *gin.Context) {
	var asc models.AppSessionContext
	reqBody, err := gc.GetRawData()
	if err != nil {
		logger.SBILog.Errorf("Get Request Body error: %+v", err)
		pd := openapi.ProblemDetailsSystemFailure(err.Error())
		gc.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		gc.JSON(http.StatusInternalServerError, pd)
		return
	}

	if err := openapi.Deserialize(&asc, reqBody, "application/json"); err != nil {
		logger.SBILog.Errorf("Deserialize Request Body error: %+v", err)
		pd := openapi.ProblemDetailsMalformedReqSyntax(err.Error())
		gc.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		gc.JSON(http.StatusBadRequest, pd)
		return
	}

	s.Processor().PostAfSessionQos(gc, gc.Param("afId"), &asc)
}

func (s *Server) apiGetAfSessionQos(gc *gin.Context) {
	s.Processor().GetAfSessionQos(gc, gc.Param("afId"), gc.Param("sessId"))
}

func (s *Server) apiPatchAfSessionQos(gc *gin.Context) {
	var ascUpdate models.AppSessionContextUpdateData
	reqBody, err := gc.GetRawData()
	if err != nil {
		logger.SBILog.Errorf("Get Request Body error: %+v", err)
		pd := openapi.ProblemDetailsSystemFailure(err.Error())
		gc.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		gc.JSON(http.StatusInternalServerError, pd)
		return
	}

	if err := openapi.Deserialize(&ascUpdate, reqBody, "application/json"); err != nil {
		logger.SBILog.Errorf("Deserialize Request Body error: %+v", err)
		pd := openapi.ProblemDetailsMalformedReqSyntax(err.Error())
		gc.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		gc.JSON(http.StatusBadRequest, pd)
		return
	}

	s.Processor().PatchAfSessionQos(gc, gc.Param("afId"), gc.Param("sessId"), &ascUpdate)
}

func (s *Server) apiDeleteAfSessionQos(gc *gin.Context) {
	s.Processor().DeleteAfSessionQos(gc, gc.Param("afId"), gc.Param("sessId"))
}
