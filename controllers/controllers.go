package controllers

import (
	"fmt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	"log"
	"main/entity"
	"main/helper"
	"os"
)

var TWILIO_ACCOUNT_SID string = os.Getenv("TWILIO_ACCOUNT_SID")
var TWILIO_AUTH_TOKEN string = os.Getenv("TWILIO_AUTH_TOKEN")
var VERIFY_SERVICE_SID string = os.Getenv("VERIFY_SERVICE_SID")
var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: TWILIO_ACCOUNT_SID,
	Password: TWILIO_AUTH_TOKEN,
})

var (
	secretkey string = "secretkeyjwt"
)

func SendOTP(to string) error {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
	}
	return err
}

func CheckOTP(to entity.Verification) {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to.Phone)
	params.SetCode(to.Code)

	resp, err := client.VerifyV2.CreateVerificationCheck(VERIFY_SERVICE_SID, params)
	if err != nil {
		log.Println("controllers", err)
	}
	if *resp.Status == "approved" {
		validToken, err := helper.GenerateJWT(to.Phone)
		if err != nil {
			log.Println("controllers 1", err)
		}

		var token helper.Token
		token.Email = to.Phone
		token.TokenString = validToken

		fmt.Println("token : ", token)
		fmt.Println("Correct!")
	} else {
		fmt.Println("Incorrect!")
	}
}
