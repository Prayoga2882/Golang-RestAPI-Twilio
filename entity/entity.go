package entity

import "time"

type User struct {
	Id         int       `json:"id"`
	Phone      string    `validate:"required" json:"phone"`
	Receiver   string    `json:"receiver"`
	Payload    string    `json:"payload"`
	VerifiedAt time.Time `json:"verifiedAt"`
	ExpiredAt  time.Time `json:"expiredAt"`
}

type Verification struct {
	Id         int       `json:"id"`
	Code       string    `validate:"required" json:"code"`
	Phone      string    `validate:"required" json:"phone"`
	Receiver   string    `validate:"required" json:"receiver"`
	Payload    string    `validate:"required" json:"payload"`
	VerifiedAt time.Time `json:"verifiedAt"`
	ExpiredAt  time.Time `json:"expiredAt"`
}

type Request struct {
	Phone      string    `validate:"required,e164" json:"phone"`
	Receiver   string    `json:"receiver"`
	Payload    string    `json:"payload"`
	VerifiedAt time.Time `json:"verifiedAt"`
	ExpiredAt  time.Time `json:"expiredAt"`
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
