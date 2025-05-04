package user_handler

import (
	"time"

	"github.com/guregu/null"
)

type UserParams struct {
	CompanyId       string      `json:"company_id"`
	UserId          string      `json:"user_id"`
	Name            string      `json:"name" validate:"required"`
	RoleID          string      `json:"role_id"`
	RoleName        string      `json:"role_name"`
	Email           string      `json:"email" validate:"required"`
	Password        string      `json:"password"`
	Mobile          string      `json:"mobile"`
	ReportingUserId null.String `json:"reporting_user_id"`
	Address         null.String `json:"address"`
	PincodeId       null.Int    `json:"pincode_id"`
	Pincode         null.String `json:"pincode"`
	CreatedAt       time.Time   `json:"created_at"`
	Status          int         `json:"status"`
}

var UserInfoRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Mobile   string `json:"mobile" validate:"required"`
	Avatar   string `json:"avatar" validate:"required"`
	GoogleID string `json:"google_id" validate:"required"`
}
