package api

import (
	"net/http"
	"server-management/service"
	"server-management/util"
	"time"

	"github.com/gin-gonic/gin"
)

type Report struct {
	report *service.ReportService
}

func NewReport(report *service.ReportService) *Report {
	return &Report{
		report: report,
	}
}

type reportRequest struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Email string `json:"email"`
}

func (r *Report) SendReport(ctx *gin.Context) {
	log := util.NewLogger()
	var request reportRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	s, err := time.Parse("02-01-2006", request.Start)
	if err != nil {
		log.Error("error when parse date: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	e, err := time.Parse("02-01-2006", request.End)
	if err != nil {
		log.Error("error when parse date: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startM := s.UnixMilli()
	endM := e.UnixMilli()
	err = r.report.SendReport(startM, endM, request.Email)
	if err != nil {
		log.Error("error when send report to email: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Info("sent report to email: " + request.Email)
	ctx.JSON(http.StatusOK, gin.H{"msg": "sent successfully"})
}
