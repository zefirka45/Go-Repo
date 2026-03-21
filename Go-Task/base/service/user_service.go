//Бизнес-логика но сейчас вызывает репозиторий 
package service

import(
	"context"
	"errors"
	"base/models"
	"base/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo:repo}
}

//Регистрация пользователя
func (s *UserService) RegisterUser(ctx context.Context, name, password, status, role, organization string) (*models.Users ,error) {
	if name == "" {
		return nil, error.New("Имя не может быть пустым")
	}
	//Подготовка сущности
	user := &models.Users{
		Name: 			name,
		Password:		password,
		Status:			"Активный",
		Role:			"Пользователь",
		Organization:	organization,

	}
	//Делегирование работы в репозиторий
	if err = s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*models.Users, error){
	return s.repo.GetByID(ctx, id)
}