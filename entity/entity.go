package entity

import "time"

type User struct {
	Id    int
	Phone string
}

type Verification struct {
	Id         int       `json:"id"`
	Code       string    `json:"code"`
	Phone      string    `json:"phone"`
	Receiver   string    `json:"receiver"`
	Payload    string    `json:"payload"`
	VerifiedAt time.Time `json:"verifiedAt"`
	ExpiredAt  time.Time `json:"expiredAt"`
}

type Request struct {
	Phone string `validate:"required,e164" json:"phone"`
}

type RequestVerification struct {
	Id    int
	Phone string `validate:"required, e164" json:"phone"`
	Code  string `validate:"required" json:"code"`
}

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
