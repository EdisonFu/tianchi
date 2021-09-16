package services

import (
	"errors"
	"tianchi/dao/cache"
	db "tianchi/dao/mysql"
	"tianchi/models"
)

func CreateUser(user *models.User) (err error) {
	if user == nil {
		return errors.New("user is nil")
	}
	_, loaded := cache.UserInfoMap.LoadOrStore(user.Username, user)
	if loaded {
		return errors.New("user is exist")
	}

	db.WriteUserChan <- user
	return
}

func Login(user *models.User) bool {
	if user == nil {
		return false
	}
	password, ok := cache.UserInfoMap.Load(user.Username)
	if !ok || password.(*models.User).Password != user.Password {
		return false
	}
	return true
}

func GetUserInfo(name string) (userResponse *models.UserResponse, err error) {
	userResponse = new(models.UserResponse)

	value, ok := cache.UserInfoMap.Load(name)
	if !ok || value == nil {
		return nil, errors.New("no user info!")
	}

	user := value.(*models.User)
	userResponse.Phone = user.Phone
	userResponse.Email = user.Email
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	return
}
