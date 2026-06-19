package repository

import (
	"errors"
	"fmt"
	"go-server/config"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) (*UserRepo, error) {
	if db == nil {
		return nil, fmt.Errorf("db connection is nil")
	}
	return &UserRepo{db: db}, nil
}
func (r *UserRepo) Create(userlogin, username, password string) error {
	hashPassword, err := HashPassword(password) //подставляем
	if err != nil {
		return err
	}
	err = r.db.Create(&config.Users{
		Userlogin: userlogin,
		Username:  username,
		Role:      "user",
		Password:  hashPassword,
		Status:    true,
	}).Error
	return err
}

func (r *UserRepo) GetUserByCredentials(userlogin, password string) (*config.Users, error) {
	var user config.Users

	err := r.db.Where("Userlogin = ? AND Status=?", userlogin, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	matched := CheckPasswordHash(password, user.Password)
	if !matched {
		return nil, errors.New("invalid password")
	}

	return &user, nil

}

// func (r *UserRepo) GetByID(id uint) (*config.Users, error) {
// 	var user config.Users
// 	err := r.db.First(&user, id).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }
