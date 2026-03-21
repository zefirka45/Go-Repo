//Выполнение CRUD операций

package repository

import (
	"contex"
	"base/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

//NewUserRepository создаем новый репозиторий 
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
//Сохраняем пользователей в бд
func (r *UserRepository) Create(ctx context.Context, user *models.Users) error {
	return r.db.WithContext(ctx).Create(user).Error
}

//GetByID ищем пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.Users, error) {
	var user models.Users
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		fmt.Println(err)
	}
	return &user,nil
}
//AutoMigrate создает таблицы
func(r *UserRepository) Migrate() error {
	return r.db.AutoMigrate(&models.Users{})
}