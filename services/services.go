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
	err := service.validate.Struct(request)
	if err != nil {
		log.Println("SERVICES CREATE")
		panic(helper.NewHandleError("MUST BE +62"))
	}

	requestClient := entity.User{
		Phone:      request.Phone,
		Receiver:   request.Receiver,
		Payload:    request.Payload,
		VerifiedAt: time.Now(),
		ExpiredAt:  time.Now().Add(3 * time.Minute),
	}

	if helper.UserExists(ctx, service.db, request.Phone) {
		panic(helper.NewHandleError("ALREADY USED GUYS"))
	}

	requestFinal, err := service.OTPrepository.Create(ctx, service.db, requestClient)
	if err != nil {
		log.Println("SERVICES CREATE 1")
		panic(helper.NewHandleError(err.Error()))
	}

	//err = controllers.SendOTP(requestClient.Phone)
	//if err != nil {
	//	fmt.Println("SERVICES CREATE 2")
	//	panic(helper.NewHandleError(err.Error()))
	//}

	return helper.UserToResponse(requestFinal), nil
}

func (service *OTPservicesImplementation) Verification(ctx context.Context, request entity.Verification) (entity.Response, error) {
	err := service.validate.Struct(request)
	if err != nil {
		fmt.Println("SERVICE VERIFICATION")
		panic(helper.NewHandleError(err.Error()))
	}

	requestClient := entity.Verification{
		Id:         request.Id,
		Code:       request.Code,
		Phone:      request.Phone,
		Receiver:   request.Receiver,
		Payload:    request.Payload,
		VerifiedAt: time.Now(),
		ExpiredAt:  time.Now().Add(3 * time.Minute),
	}

	userData := entity.User{
		Id:         request.Id,
		Phone:      request.Phone,
		Payload:    request.Payload,
		VerifiedAt: time.Now(),
		ExpiredAt:  time.Now().Add(3 * time.Minute),
	}

	today := time.Now()
	expiredAt := userData.ExpiredAt
	if today == expiredAt {
		panic(helper.NewHandleError("EXPIRED"))
	}
	err = controllers.CheckOTP(requestClient)
	if err != nil {
		fmt.Println("SERVICE VERIFICATION 2", err)
		panic(helper.NewHandleError(err.Error()))
	}

	_, err = service.OTPrepository.Verification(ctx, service.db, requestClient)
	if err != nil {
		fmt.Println("SERVICE VERIFICATION 1")
		panic(helper.NewHandleError(err.Error()))
	}

	return helper.RequestVerificationToResponse(requestClient), nil
}
