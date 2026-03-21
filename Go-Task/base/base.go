//Сборка модулей
package base
import (
	"contex"
	"fmt"
	"log"

	"base/repository"
	"base/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Base() {
	//Инициализация БД
	db, err := gorm.Open(sqlite.Open("unit.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	ctx := context.Background()
	
	//Сборка
	userRepo := repository.NewUserRepository(db)
	
	//Миграция таблиц
	if err := userRepo.Migrate();err != nil {
		log.Fatalf("Ошибка миграции: %v",err)
	}

	userService := service.NewUserServie(userRepo)

	fmt.Println("Запуск приложения")

	newUser, err := userService.RegisterUser(ctx, "", "","","")
	if err != nil {
		fmt.Printf("Ошибка создания: %v\n",err)
	} else {
		fmt.Printf("Пользователь создан: %+v\n",newUser)
	}

	foundUser, err := userService.GetUser(ctx, ?)
	if err != nil{
		fmt.Printf("Пользователь не найден: %v\n", err)
	} else {
		fmt.Printf("Пользователь найден: %+v\n", foundUser)
	}
}