package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"main/entity"
	"main/helper"
)

type OTPrepository interface {
	Create(ctx context.Context, db *sql.DB, user entity.User) (entity.User, error)
	Verification(ctx context.Context, db *sql.DB, request entity.Verification) (entity.Verification, error)
}

type OTPrepositoryImplementation struct{}

func NewOTPrepositoryImplementation() *OTPrepositoryImplementation {
	return &OTPrepositoryImplementation{}
}

func (otp *OTPrepositoryImplementation) Create(ctx context.Context, db *sql.DB, user entity.User) (entity.User, error) {
	if helper.UserExists(db, user.Phone) {
		return user, errors.New("phone already used")
	}

	sql := "INSERT INTO users(phone) VALUES (?)"
	execContext, err := db.ExecContext(ctx, sql, user.Phone)
	helper.HandlePanic(err)

	id, err := execContext.LastInsertId()
	helper.HandlePanic(err)
	user.Id = int(id)

	return user, nil
}

func (otp *OTPrepositoryImplementation) Verification(ctx context.Context, db *sql.DB, verified entity.Verification) (entity.Verification, error) {
	sql := "INSERT INTO verification (code, phone, receiver, payload) VALUES (?, ?, ?, ?)"
	execContext, err := db.ExecContext(ctx, sql, verified.Code, verified.Phone, verified.Receiver, verified.Payload)
	if err != nil {
		log.Println("REPOSITORY", err)
	}
	id, err := execContext.LastInsertId()
	if err != nil {
		log.Println("REPOSITORY 1", err)
	}
	verified.Id = int(id)
	return verified, nil
}
