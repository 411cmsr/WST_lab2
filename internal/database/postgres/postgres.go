package postgres

import (
	"WST_lab1_server_new1/config"
	"WST_lab1_server_new1/internal/database"
	"WST_lab1_server_new1/internal/logging"
	"WST_lab1_server_new1/internal/models"

	"errors"
	"strconv"
	"strings"

	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*

 */

type Storage struct {
	DB               *gorm.DB
	PersonRepository *PersonRepository
}

type PersonRepository struct {
	DB *gorm.DB
}

/*
Инициализация
*/
func Init() (*Storage, error) {
	logging.InitializeLogger()
	var err error
	//Уровень логирования из файла конфигурации
	var logLevel logger.LogLevel
	switch config.GeneralServerSetting.LogLevel {
	case "fatal":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info", "debug":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}
	//Строка подключения к базе данных
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.DatabaseSetting.Host,
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Name,
		config.DatabaseSetting.Port,
		config.DatabaseSetting.SSLMode)
	//Подключаемся к базе данных
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	//Выводим при удачном подключении
	logging.Logger.Info("Database connection established successfully.")
	//Миграция базы данных
	db := conn
	err = db.AutoMigrate(&models.Person{})
	if err != nil {
		log.Fatalf("error creating table: %v", err)
		return nil, fmt.Errorf("error creating table: %v", err)
	}
	logging.Logger.Info("Migration completed successfully.")
	//Удаляем таблицу
	db.Exec("DELETE FROM people")
	//Заполняем таблицу из фаила конфигурации
	result := db.Create(&config.GeneralServerSetting.DataSet)
	if result.Error != nil {
		log.Fatalf("error creating table: %v", result.Error)
	}
	//Выводим при удачном заполнении таблицы
	logging.Logger.Info("Database updated successfully.")

	/*
		//Debug: Запрос к базе и вывод всех данных
	*/
	var results []models.Person
	if err := db.Find(&results).Error; err != nil {
		log.Fatalf("query failed: %v", err)
	}
	for _, record := range results {
		fmt.Println(record)

	}
	fmt.Println("database content in quantity:", len(results), "\n id max:", results[len(results)-1].ID, "id min:", results[0].ID)
	/*
		----
	*/
	//Возвращаем указатель

	personRepo := &PersonRepository{DB: db}
	return &Storage{
		DB:               db,
		PersonRepository: personRepo,
	}, nil

}

/*
//
Метод поиска в базе данных по запросу
//
*/
func (pr *PersonRepository) SearchPerson(searchString string) ([]models.Person, error) {
	var persons []models.Person
	query := pr.DB.Model(&models.Person{})
	//Удаляем пробелы из строки поиска
	searchString = strings.TrimSpace(searchString)
	// Проверяем строка является числом, если число ищем по возрасту
	if age, err := strconv.Atoi(searchString); err == nil {
		query = query.Where("age = ?", age)
	} else {
		//Если строка не может быть конвертирована в число ищем по строковым полям
		query = query.Where("name LIKE ? OR surname LIKE ? OR email LIKE ? OR telephone LIKE ?",
			"%"+searchString+"%", "%"+searchString+"%", "%"+searchString+"%", "%"+searchString+"%")
	}
	//Выполняем запрос и сохраняем результат в структуру
	if err := query.Find(&persons).Error; err != nil {
		return nil, err
	}
	//Возвращаем результат
	return persons, nil
}

/*
Метод добавления новых данных
*/
func (pr *PersonRepository) AddPerson(person *models.Person) (uint, error) {
	//Проверяем наличие записи с таким же email
	if _, err := pr.CheckPersonByEmail(person.Email, 0); err == nil {
		return 0, database.ErrEmailExists
	}
	//Создаем запись в базе данных
	if err := pr.DB.Create(person).Error; err != nil {
		return 0, err
	}
	//Возвращаем id созданной записи
	return person.ID, nil
}

/*
Метод получения данных по id
*/
func (pr *PersonRepository) GetPerson(id uint) (*models.Person, error) {
	var person models.Person
	//Выполняем запрос к базе данных для получения записи по id
	err := pr.DB.First(&person, id).Error
	if err != nil {
		//Возвращаем ошибку при выполнении запроса к базе данных
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrPersonNotFound
		}
		return nil, err
	}
	//Возвращаем результат
	return &person, nil
}

/*
Метод обновления данных по id
*/
func (pr *PersonRepository) UpdatePerson(person *models.Person) error {
	//Выполняем запрос к базе данных для обновления записи
	if _, err := pr.CheckPersonByEmail(person.Email, person.ID); err == nil {
		return database.ErrEmailExists
	}
	//Выполняем запрос к базе данных для обновления записи
	result := pr.DB.Model(&models.Person{}).Where("id = ?", person.ID).Updates(models.Person{
		Name:      person.Name,
		Surname:   person.Surname,
		Age:       person.Age,
		Email:     person.Email,
		Telephone: person.Telephone,
	})

	if result.Error != nil {
		//Возвращаем ошибку при выполнении запроса к базе данных
		return result.Error
	}

	if result.RowsAffected == 0 {
		//Возвращаем ошибку если запись не найдена для обновления
		return database.ErrPersonNotFound
	}
	//Возвращаем ничего при успехе
	return nil
}

/*
Метод удаления данных по id
*/
func (pr *PersonRepository) DeletePerson(request *models.DeletePersonRequest) error {
	if err := pr.DB.Delete(&models.Person{}, request.ID).Error; err != nil {
		return err
	}
	return nil
}

/*
Метод получения всех данных
*/
func (pr *PersonRepository) GetAllPersons() ([]models.Person, error) {
	var persons []models.Person
	//Выполняем запрос к базе данных для получения всех записей
	err := pr.DB.Find(&persons).Error
	if err != nil {
		return nil, err
	}
	//Возвращаем результат
	return persons, nil
}

/*
Метод проверки наличия записи по email
*/
func (pr *PersonRepository) CheckPersonByEmail(email string, excludeId uint) (*models.Person, error) {
	var person models.Person
	// Выполняем запрос к базе данных для поиска по email
	if err := pr.DB.Where("email = ? AND id != ?", email, excludeId).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//Возвращаем кастомную ошибку (Запись не найдена)
			return nil, database.ErrPersonNotFound
		}

		//Возвращаем ошибку
		return nil, err
	}
	//Возвращаем запись

	return &person, nil
}

/*
Метод проверки наличия записи по id
*/
func (pr *PersonRepository) CheckPersonByID(id uint) (bool, error) {
	var person models.Person
	//Выполняем запрос к базе данных для поиска по id
	result := pr.DB.First(&person, id)
	if result.Error != nil {
		//Проверяем наличие записи по id
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			//fmt.Println("Record not found")
			return false, database.ErrPersonNotFound
		} else {

			//fmt.Println("Error when executing the request:", result.Error)
			return false, result.Error
		}
	} else {
		fmt.Println("The record was found with CheckPersonByIDHandler:", person)
		return true, nil
	}
	//Возвращаем false при успехе
	return false, nil
}
