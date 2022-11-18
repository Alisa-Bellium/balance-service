package services

import (
	"balance-service/repositories"
	"errors"
	_ "github.com/lib/pq"
)

var ErrUserNotExists = errors.New("user does not exists")

func GetUserBalance(userID int64) (*repositories.User, int64, error) {
	user, err := repositories.GetUser(nil, userID)

	if err != nil {
		return nil, 0, err
	}
	if user == nil {
		return nil, 0, ErrUserNotExists
	}

	reserved, err := repositories.GetUserReservedAmount(nil, userID)
	if err != nil {
		return nil, 0, err
	}

	return user, reserved, nil
}
