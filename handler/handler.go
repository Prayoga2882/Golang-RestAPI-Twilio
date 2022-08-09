package handler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"main/entity"
	"main/helper"
	"main/services"
	"net/http"
)

type OTPhandler interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Verification(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type OTPhandlerImplementation struct {
	OTPservice services.OTPservices
}

func NewOTPhandlerImplementation(OTPservice services.OTPservices) *OTPhandlerImplementation {
	return &OTPhandlerImplementation{OTPservice: OTPservice}
}

func (handler *OTPhandlerImplementation) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	requestBody := entity.Request{}
	helper.ReadFromRequestBody(request, &requestBody)

	_, err := handler.OTPservice.Create(request.Context(), requestBody)
	if err != nil {
		var err helper.Error
		err = helper.SetError(err, "ADA YANG SALAH")
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(err)
		return
	}

	responseBody := entity.Response{
		Code:   200,
		Status: "Successfully",
		Data:   requestBody.Phone,
	}
	helper.WriteToResponseBody(writer, responseBody)
}

func (handler *OTPhandlerImplementation) Verification(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	requestBody := entity.RequestVerification{}
	helper.ReadFromRequestBody(request, &requestBody)

	requestBodyResponse := helper.RequestVerificationToVerification(requestBody)
	_, err := handler.OTPservice.Verification(request.Context(), requestBodyResponse)
	if err != nil {
		log.Println("HANDLER VERIFICATION", err)
	}

	responseBody := entity.Response{
		Code:   200,
		Status: "Successfully",
		Data:   requestBody.Phone,
	}

	helper.WriteToResponseBody(writer, responseBody)

}
