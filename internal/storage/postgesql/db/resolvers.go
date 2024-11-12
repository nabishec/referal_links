package db

import (
	"fmt"

	"time"

	"github.com/nabishec/referal_links/internal/models"
	"github.com/rs/zerolog/log"
)

func (r *Database) AddUser(user *models.UserInfo) error {
	const op = "internal.storage.postgresql.db.AddUser()"

	log.Debug().Str("email", user.Email).Msg("Request to the database to add new user")

	if _, err := r.FoundUserId(user); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	user.Date = time.Now().Format("2006-01-02")
	_, err := r.DB.Exec("INSERT INTO users (name, password, email, date) VALUES ($1,$2,$3,$4)", user.Name, user.Password, user.Email, user.Date)

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msg("New user added to db")

	return nil
}

func (r *Database) FoundUserId(user *models.UserInfo) (int64, error) {
	op := "internal.storage.postgresql.db.FoundUserId()"

	log.Debug().Str("email", user.Email).Msg("Found user id in db")

	var userId int64
	err := r.DB.QueryRow("SELECT id FROM users WHERE email = $1", user.Email).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Int64("user", userId).Msg("found")

	return userId, nil
}

func (r *Database) AddReferral(referralUser *models.UserInfo, email string) error {
	op := "internal.storage.postgresql.db.AddReferral()"

	log.Debug().Str("email", email).Msg("Request to the database to add a new referral")

	//TODO transaction
	if err := r.AddUser(referralUser); err != nil {
		return fmt.Errorf("%s:%w(%s)", op, err, "can't add in users table")
	}

	referrer := &models.UserInfo{Email: email}
	referrerId, err := r.FoundUserId(referrer)
	if err != nil {
		return fmt.Errorf("%s:%w(%s)", op, err, "referrer not found")
	}

	referral := &models.Referral{
		ReferrerId: referrerId,
		Name:       referralUser.Name,
		Email:      referralUser.Email,
		Date:       time.Now().Format("2006-01-02"),
	}

	_, err = r.DB.Exec("INSERT INTO referrals (referrer_id, referral_name, referral_email, date) VALUES ($1,$2,$3,$4)", referral.ReferrerId, referral.Name, referral.Email, referral.Date)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msg("New referral added to db")

	return nil
}

func (r *Database) FoundReferrals(referrerId int64) ([]*models.ReferralResponse, error) {
	op := "internal.storage.postgresql.db.FoundReferrals()"

	log.Debug().Int64("referrer", referrerId).Msg("Found referrals")

	var referrals []*models.ReferralResponse
	err := r.DB.Select(&referrals, "SELECT referral_name, date FROM referrals WHERE referrer_id = $1", referrerId)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msg("Referrals found")

	return referrals, nil
}
