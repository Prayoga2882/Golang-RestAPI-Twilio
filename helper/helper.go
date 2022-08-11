package helper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"log"
	"main/entity"
	"net/http"
	"time"
)

const (
	StatusMessageOK string = "OK"

	// StatusMessageBadRequest is custom status message for bad request
	StatusMessageBadRequest string = "Bad Request"

	// StatusMessageInternalServerError is custom status message for unknown error / internal server error
	StatusMessageInternalServerError string = "Internal Error"

	// StatusMessageNotFound is custom status message for data not found
	StatusMessageNotFound string = "Not Found"

	// StatusMessageForbidden is custom status message for forbidden
	StatusMessageForbidden string = "Forbidden"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ErrorHandler struct {
	Err           error
	Status        string
	MessageStatus string
	HTTPStatus    int
}

func NewErrorHandler(err error, status string, messageStatus string, HTTPStatus int) *ErrorHandler {
	return &ErrorHandler{Err: err, Status: status, MessageStatus: messageStatus, HTTPStatus: HTTPStatus}
}

func ErrBadRequest(err error, message string) *ErrorHandler {
	if len(message) <= 0 || message == "" {
		message = StatusMessageBadRequest
	}
	return &ErrorHandler{
		Err:           err,
		Status:        StatusMessageBadRequest,
		MessageStatus: message,
		HTTPStatus:    404,
	}
}

var (
	secretkey string = "secretkeyjwt"
)

type Token struct {
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

func NotFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(Error)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := entity.Response{
			Code:   404,
			Status: "NOT FOUND",
			Data:   exception,
		}
		WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func ValidationErrors(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := entity.Response{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   exception,
		}

		WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func InternalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	webResponse := entity.Response{
		Code:   500,
		Status: "INTERNAL SERVER ERROR",
		Data:   nil,
	}
	WriteToResponseBody(writer, webResponse)

}

func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}

func GenerateJWT(phone string) (string, error) {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["phone"] = phone
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("something went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func HandlePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func RequestToUser(response entity.Request) entity.User {
	return entity.User{
		Phone: response.Phone,
	}
}

func RequestVerificationToVerification(response entity.RequestVerification) entity.Verification {
	return entity.Verification{
		Id:         response.Id,
		Code:       response.Code,
		Phone:      response.Phone,
		Receiver:   response.Phone,
		Payload:    response.Phone,
		VerifiedAt: time.Now(),
		ExpiredAt:  time.Now(),
	}
}

func RequestVerificationToResponse(response entity.Verification) entity.Response {
	return entity.Response{
		Code:   200,
		Status: "Succesfully",
		Data:   response,
	}
}

func UserToResponse(response entity.User) entity.Response {
	return entity.Response{
		Code:   200,
		Status: "Successfully",
		Data:   response.Phone,
	}
}

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	HandlePanic(err)
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	HandlePanic(err)
}

func UserExists(db *sql.DB, phone string) bool {
	sqlStmt := `SELECT phone FROM users WHERE phone = ?`
	err := db.QueryRow(sqlStmt, phone).Scan(&phone)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}
		return false
	}
	return true
}
