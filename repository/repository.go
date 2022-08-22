package repository

import (
	"context"
	"database/sql"
	"fmt"
	"main/entity"
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
	sql := "INSERT INTO users (phone, receiver, payload, verified_at, expired_at) VALUES (?, ?, ?, ?, ?)"
	execContext, err := db.ExecContext(ctx, sql, user.Phone, user.Receiver, user.Payload, user.VerifiedAt, user.ExpiredAt)
	if err != nil {
		fmt.Println("REPOSITORY")
		return user, err
	}
	id, err := execContext.LastInsertId()
	if err != nil {
		fmt.Println("REPOSITORY 1")
		return user, err
	}
	user.Id = int(id)
	return user, nil
}

func (otp *OTPrepositoryImplementation) Verification(ctx context.Context, db *sql.DB, verified entity.Verification) (entity.Verification, error) {
	sql := "INSERT INTO verification (code, phone, verified_at, expired_at) VALUES (?, ?, ?, ?)"
	execContext, err := db.ExecContext(ctx, sql, verified.Code, verified.Phone, verified.VerifiedAt, verified.ExpiredAt)
	if err != nil {
		fmt.Println("REPOSITORY")
		return verified, err
	}
	id, err := execContext.LastInsertId()
	if err != nil {
		fmt.Println("REPOSITORY 1")
		return verified, err
	}
	verified.Id = int(id)
	return verified, nil
}
