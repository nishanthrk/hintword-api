package user_handler

import (
	"github.com/guregu/null"
	"time"
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
