package models

import (
	"database/sql"
	"strconv"
	"time"
)

type User struct {
	ID              int            `db:"id"`
	Name            string         `db:"name"`
	Age             int            `db:"age"`
	Phone           string         `db:"phone"`
	Email           string         `db:"email"`
	Status          int            `db:"status"`
	BlockReason     sql.NullString `db:"block_reason"`       // notice: lowercase in db
	BlockReasonCode int            `db:"blockr_reason_code"` // notice: lowercase in db
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

type UserData struct {
	Users   []User
	Filter  UserRequest
	Reasons map[int]string
}

type UserRequest struct {
	UserId             int64    `json:"userId"`     // Phone number or user ID
	ReasonCode         string   `json:"reasonCode"` // Reason code (1, 2, 9999)
	Reason             string   `json:"reason"`     // Custom reason for "Other"
	BlockReasonCodeInt *int     `json:"-"`
	ReasonCodeInt      int      `json:"-"`
	Name               string   `json:"name"`
	Phone              []string `josn:"-"`
	Email              string   `json:"email"`
	Status             string   `json:"status"`
	BlockReason        string   `json:"blockReason"`
	StatusInt          *int     `json:"-"`
	PhoneStr           string   `json:"phone"`
}

type AdminCred struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Admin struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
}

func (req *UserRequest) HandleIntFields() error {
	if req.ReasonCode != "" {
		code, err := strconv.Atoi(req.ReasonCode)
		if err != nil {
			return err
		}
		req.ReasonCodeInt = code
	}

	if req.BlockReason != "" {
		reason, err := strconv.Atoi(req.BlockReason)
		if err != nil {
			return err
		}
		req.BlockReasonCodeInt = &reason
	}

	if req.Status != "" {
		status, err := strconv.Atoi(req.Status)
		if err != nil {
			return err
		}
		req.StatusInt = &status
	}

	return nil
}
