package helper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"main/entity"
	"net/http"
	"time"
)

var (
	secretkey string = "secretkeyjwt"
)

type Token struct {
	TokenString string `json:"token"`
}

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

type HandleError struct {
	Error string
}

func NewHandleError(error string) *HandleError {
	return &HandleError{Error: error}
}

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {

	ValidationError(writer, request, err)
}

func BadRequest(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)

	webResponse := entity.Response{
		Code:   400,
		Status: "BAD REQUEST",
		Data:   err,
	}
	WriteToResponseBody(writer, webResponse)
}

func ValidationError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnprocessableEntity)

	webResponse := entity.Response{
		Code:   422,
		Status: "UNPROCESSABLE",
		Data:   err,
	}
	WriteToResponseBody(writer, webResponse)
}

func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}

func GenerateJWT() (string, error) {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
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

func UserExists(ctx context.Context, db *sql.DB, phone string) bool {
	sqlStmt := `SELECT phone FROM users WHERE phone = ?`
	err := db.QueryRowContext(ctx, sqlStmt, phone).Scan(&phone)
	if err != nil {
		return false
	}
	return true
}
