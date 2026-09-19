package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hospital-backend/internal/model"
	"hospital-backend/internal/service"
)

type PatientHandler struct {
	patientService *service.PatientService
}

func NewPatientHandler(patientService *service.PatientService) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}

func (handler *PatientHandler) SearchPatients(context *gin.Context) {
	hospitalIDValue, exists := context.Get("hospital_id")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	hospitalID, ok := hospitalIDValue.(int)
	if !ok || hospitalID <= 0 {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid hospital scope"})
		return
	}

	var request model.PatientSearchRequest
	if err := context.ShouldBindQuery(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid search query",
		})
		return
	}

	patients, err := handler.patientService.Search(
		context.Request.Context(),
		hospitalID,
		request,
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"count": len(patients),
		"data":  patients,
	})
}
