package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"main/database"
	"main/handler"
	"main/helper"
	"main/repository"
	"main/services"
	"net/http"
)

func main() {
	db := database.NewDB()
	validate := validator.New()
	repositoryOTP := repository.NewOTPrepositoryImplementation()
	servicesOTP := services.NewOTPserviceImplementation(repositoryOTP, db, validate)
	handlerOTP := handler.NewOTPhandlerImplementation(servicesOTP)

	router := httprouter.New()

	router.POST("/api/register", handlerOTP.Create)
	router.POST("/api/verification", handlerOTP.Verification)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}
	fmt.Println("Application is running...")
	err := server.ListenAndServe()
	helper.HandlePanic(err)
}
