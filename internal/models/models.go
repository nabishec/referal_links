package models

type UserInfo struct {
	Name     string `json:"name" validate:"required" db:"name"`
	Email    string `json:"email" validate:"required,email" db:"email"`
	Password string `json:"password" validate:"required" db:"password"`
	Date     string `db:"date"`
}

type ReferralInfo struct {
	ReferralCode string `json:"referral_cade" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
}

type Referral struct {
	ReferrerId int64  `db:"referrer_id"`
	Name       string `db:"referral_name"`
	Email      string `db:"referral_email"`
	Date       string `db:"date"`
}

type ReferralResponse struct {
	Name string `json:"referral_name" db:"referral_name"`
	Date string `json:"registration_date" db:"date"`
}
