package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"main/controllers"
	"main/entity"
	"main/helper"
	"main/repository"
	"time"
)

type OTPservices interface {
	Create(ctx context.Context, request entity.Request) (entity.Response, error)
	Verification(ctx context.Context, request entity.Verification) (entity.Response, error)
}

type OTPservicesImplementation struct {
	OTPrepository repository.OTPrepository
	db            *sql.DB
	validate      *validator.Validate
}

func NewOTPserviceImplementation(OTPrepository repository.OTPrepository, db *sql.DB, validate *validator.Validate) *OTPservicesImplementation {
	return &OTPservicesImplementation{OTPrepository: OTPrepository, db: db, validate: validate}
}

func (service *OTPservicesImplementation) Create(ctx context.Context, request entity.Request) (entity.Response, error) {

	requestClient := entity.User{
		Phone:      request.Phone,
		Receiver:   request.Receiver,
		Payload:    request.Payload,
		VerifiedAt: time.Now(),
		ExpiredAt:  time.Now().Add(3 * time.Minute),
	}

	err := service.validate.Struct(request)
	if err != nil {
		log.Println("SERVICES CREATE")
		panic(helper.NewHandleError("CREATE MUST BE +62"))
	}

	if requestClient.Receiver == "" {
		panic(helper.NewHandleError("RECEIVER CANNOT BE EMPTY"))
	}

	if requestClient.Payload == "" {
		panic(helper.NewHandleError("PAYLOAD CANNOT BE EMPTY"))
	}

	if helper.UserExists(ctx, service.db, requestClient.Phone) {
		panic(helper.NewHandleError("PHONE ALREADY USED"))
	}

	requestFinal, err := service.OTPrepository.Create(ctx, service.db, requestClient)
	if err != nil {
		log.Println("SERVICES CREATE 1")
		panic(helper.NewHandleError(err.Error()))
	}

	err = controllers.SendOTP(requestClient.Phone)
	if err != nil {
		fmt.Println("SERVICES CREATE 2")
		panic(helper.NewHandleError(err.Error()))
	}

	return helper.UserToResponse(requestFinal), nil
}

func (service *OTPservicesImplementation) Verification(ctx context.Context, request entity.Verification) (entity.Response, error) {

	requestClient := entity.Verification{
		Id:         request.Id,
		Code:       request.Code,
		Phone:      request.Phone,
		VerifiedAt: time.Now(),
		ExpiredAt:  time.Now().Add(3 * time.Minute),
	}

	err := service.validate.Struct(request)
	if err != nil {
		fmt.Println("SERVICE VERIFICATION")
		panic(helper.NewHandleError("PHONE MUST BE +62"))
	}

	err = controllers.CheckOTP(requestClient)
	if err != nil {
		fmt.Println("SERVICE VERIFICATION 1", err)
		panic(helper.NewHandleError(err.Error()))
	}

	_, err = service.OTPrepository.Verification(ctx, service.db, requestClient)
	if err != nil {
		fmt.Println("SERVICE VERIFICATION 2")
		panic(helper.NewHandleError(err.Error()))
	}

	return helper.RequestVerificationToResponse(requestClient), nil
}
