package services

import (
	"context"
	"database/sql"
	"errors"
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
	validate := entity.Response{}
	err := service.validate.Struct(request)
	if err != nil {
		fmt.Println("SERVICES CREATE")
		return validate, err
	}

	requestClient := entity.Request{
		Phone: request.Phone,
	}
	requestClintFinal := helper.RequestToUser(requestClient)
	hasil, err := service.OTPrepository.GetUserByPhone(ctx, service.db, requestClintFinal)
	helper.HandlePanic(err)
	if hasil.Id > 0 {
		log.Println("UDAH ADA COK !!")
		return validate, errors.New("dah ada nih")
	}
	requestFinal := helper.RequestToUser(requestClient)

	requestFinal, err = service.OTPrepository.Create(ctx, service.db, requestFinal)
	if err != nil {
		fmt.Println("SERVICES CREATE 1")
		return validate, err
	}

	//err = controllers.SendOTP(requestClient.Phone)
	//if err != nil {
	//	fmt.Println("SERVICES CREATE 2")
	//	return validate, err
	//}

	return helper.UserToResponse(hasil), nil
}

func (service *OTPservicesImplementation) Verification(ctx context.Context, request entity.Verification) (entity.Response, error) {
	validate := entity.Response{}
	err := service.validate.Struct(request)
	if err != nil {
		fmt.Println("SERVICE VERIFICATION")
		return validate, err
	}

	requestClient := entity.Verification{
		Id:         request.Id,
		Code:       request.Code,
		Phone:      request.Phone,
		Receiver:   request.Receiver,
		Payload:    request.Payload,
		VerifiedAt: time.Now(),
		ExpiredAt:  time.Now(),
	}

	err = controllers.CheckOTP(requestClient)
	if err != nil {
		fmt.Println("SERVICE VERIFICATION 2", err)
		return validate, err
	}

	_, err = service.OTPrepository.Verification(ctx, service.db, requestClient)
	if err != nil {
		fmt.Println("SERVICE VERIFICATION 1")
		return validate, err
	}

	return helper.RequestVerificationToResponse(requestClient), nil
}
