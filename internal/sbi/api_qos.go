package sbi

import (
	"net/http"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/metrics/sbi"
	"github.com/gin-gonic/gin"
)

func (s *Server) getAsSessionQosRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/:scsAsId/subscriptions",
			APIFunc: s.apiListAsSessionQosSubs,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/:scsAsId/subscriptions",
			APIFunc: s.apiPostAsSessionQosSub,
		},
		{
			Method:  http.MethodGet,
			Pattern: "/:scsAsId/subscriptions/:subscriptionId",
			APIFunc: s.apiGetAsSessionQosSub,
		},
		{
			Method:  http.MethodPut,
			Pattern: "/:scsAsId/subscriptions/:subscriptionId",
			APIFunc: s.apiPutAsSessionQosSub,
		},
		{
			Method:  http.MethodPatch,
			Pattern: "/:scsAsId/subscriptions/:subscriptionId",
			APIFunc: s.apiPatchAsSessionQosSub,
		},
		{
			Method:  http.MethodDelete,
			Pattern: "/:scsAsId/subscriptions/:subscriptionId",
			APIFunc: s.apiDeleteAsSessionQosSub,
		},
	}
}

func (s *Server) apiListAsSessionQosSubs(gc *gin.Context) {
	s.Processor().ListAsSessionQosSubs(gc, gc.Param("scsAsId"))
}

func (s *Server) apiPostAsSessionQosSub(gc *gin.Context) {
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

	s.Processor().PostAsSessionQosSub(gc, gc.Param("scsAsId"), &asc)
}

func (s *Server) apiGetAsSessionQosSub(gc *gin.Context) {
	s.Processor().GetAsSessionQosSub(gc, gc.Param("scsAsId"), gc.Param("subscriptionId"))
}

func (s *Server) apiPutAsSessionQosSub(gc *gin.Context) {
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

	s.Processor().PutAsSessionQosSub(gc, gc.Param("scsAsId"), gc.Param("subscriptionId"), &ascUpdate)
}

func (s *Server) apiPatchAsSessionQosSub(gc *gin.Context) {
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

	s.Processor().PatchAsSessionQosSub(gc, gc.Param("scsAsId"), gc.Param("subscriptionId"), &ascUpdate)
}

func (s *Server) apiDeleteAsSessionQosSub(gc *gin.Context) {
	s.Processor().DeleteAsSessionQosSub(gc, gc.Param("scsAsId"), gc.Param("subscriptionId"))
}
