package db

import (
	"fmt"
	"time"

	"github.com/nabishec/referal_links/internal/models"
)

func (r *Database) AddUser(user *models.UserInfo) error {
	const op = "internal.storage.postgresql.db.AddUser()"

	if _, err := r.FoundUserId(user); err != nil {
		return fmt.Errorf("func:%s error:%w", op, err)
	}

	user.Date = time.Now().Format("2006-01-02")
	_, err := r.DB.Exec("INSERT INTO users (name, password, email, date) VALUES ($1,$2,$3,$4)", user.Name, user.Password, user.Email, user.Date)

	if err != nil {
		return fmt.Errorf("func:%s  error:%w", op, err)
	}

	return nil
}

func (r *Database) FoundUserId(user *models.UserInfo) (int64, error) {
	op := "internal.storage.postgresql.db.FoundUserId()"

	var userId int64
	err := r.DB.QueryRow("SELECT id FROM users WHERE email = $1", user.Email).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("func:%s  error:%w", op, err)
	}
	return userId, nil
}

func (r *Database) AddReferral(referralUser *models.UserInfo, email string) error {
	op := "internal.storage.postgresql.db.AddReferral()"
	//TODO transaction
	if err := r.AddUser(referralUser); err != nil {
		return fmt.Errorf("func:%s error:%w(%s)", op, err, "can't add in users table")
	}

	referrer := &models.UserInfo{Email: email}
	referrerId, err := r.FoundUserId(referrer)
	if err != nil {
		return fmt.Errorf("func:%s error:%w(%s)", op, err, "referrer not found")
	}

	referral := &models.Referral{
		ReferrerId: referrerId,
		Name:       referralUser.Name,
		Email:      referralUser.Email,
		Date:       time.Now().Format("2006-01-02"),
	}

	_, err = r.DB.Exec("INSERT INTO referrals (referrer_id, referral_name, referral_email, date) VALUES ($1,$2,$3,$4)", referral.ReferrerId, referral.Name, referral.Email, referral.Date)
	if err != nil {
		return fmt.Errorf("func:%s error:%w", op, err)
	}

	return nil
}

func (r *Database) FoundReferrals(referrerId int64) ([]*models.ReferralResponse, error) {
	op := "internal.storage.postgresql.db.FoundReferrals()"

	var referrals []*models.ReferralResponse
	err := r.DB.Select(&referrals, "SELECT referral_name, date FROM referrals WHERE referrer_id = $1", referrerId)
	if err != nil {
		return nil, fmt.Errorf("func:%s error:%w", op, err)
	}
	return referrals, nil
}
