package handlers

import (
	"WST_lab1_server_new1/internal/database"
	"WST_lab1_server_new1/internal/database/postgres"
	"WST_lab1_server_new1/internal/logging"
	"WST_lab1_server_new1/internal/models"
	"bytes"
	
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
Структура обработчика для разделения логики обработки запросов от доступа к данным
*/
type StorageHandler struct {
	Storage *postgres.Storage
}

// Обработчик SOAP запросов
func (sh *StorageHandler) SOAPHandler(c *gin.Context) {

	var envelope models.Envelope

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading request body")
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := xml.Unmarshal(body, &envelope); err != nil {
		fmt.Println("Error decoding XML:", err)
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	fmt.Printf("Decoded Envelope: %+v\n", envelope)

	switch {
	case envelope.Body.AddPerson != nil:
		sh.addPersonHandler(c, envelope.Body.AddPerson)
	case envelope.Body.DeletePerson != nil:
		sh.deletePersonHandler(c, envelope.Body.DeletePerson)
	case envelope.Body.UpdatePerson != nil:
		sh.updatePersonHandler(c, envelope.Body.UpdatePerson)
	case envelope.Body.GetPerson != nil:
		sh.getPersonHandler(c, envelope.Body.GetPerson)
	case envelope.Body.GetAllPersons != nil:
		sh.getAllPersonsHandler(c)
	case envelope.Body.SearchPerson != nil:
		sh.searchPersonHandler(c, envelope.Body.SearchPerson)
	default:
		fmt.Println("Unsupported action")
		c.String(http.StatusBadRequest, "Unsupported action")
		return
	}
}

// Метод добавления новой записи в базу данных
func (h *StorageHandler) addPersonHandler(c *gin.Context, request *models.AddPersonRequest) {

	person := models.Person{
		Name:      request.Name,
		Surname:   request.Surname,
		Age:       request.Age,
		Email:     request.Email,
		Telephone: request.Telephone,
	}


	// Добавляем person в базу данных
	id, err := h.Storage.PersonRepository.AddPerson(&person)
	if err != nil {
		if errors.Is(err, database.ErrEmailExists) {
			return
		}
		fmt.Printf("Error adding person: %v\n", err)

		// Формируем SOAP Fault для ошибки добавления
		
		return
	}
	fmt.Printf("Person added with ID: %d\n", id)

	response := models.AddPersonResponse{
		ID: id,
	}

	// Возвращаем успешный ответ в формате XML
	fmt.Printf("Response: %+v\n", response)
	c.XML(http.StatusOK, response)
}

// Метод обновления записи в базе данных
func (h *StorageHandler) updatePersonHandler(c *gin.Context, request *models.UpdatePersonRequest) {
	// Проверяем, существует ли запись с данным ID
	checkByID, err := h.Storage.PersonRepository.CheckPersonByID(uint(request.ID))
	if !checkByID {
		return
	}
	if err != nil {
		return
	}

	// Создаем объект типа Person на основе запроса
	person := models.Person{
		ID:        uint(request.ID),
		Name:      request.Name,
		Surname:   request.Surname,
		Age:       request.Age,
		Email:     request.Email,
		Telephone: request.Telephone,
	}

	// Обновляем информацию о человеке в базе данных
	err = h.Storage.PersonRepository.UpdatePerson(&person)
	if err != nil {
		// Проверяем, существует ли запись с данным Email кроме обновляемой
		if errors.Is(err, database.ErrEmailExists) {
		
			return
		}
		logging.Logger.Error("Error updating person with ID", zap.Uint("ID", uint(request.ID)), zap.Error(err))

		
		return
	}
	fmt.Println("Person updated with ID:")
	logging.Logger.Info("Successfully updated person with ID", zap.Uint("ID", uint(request.ID)))

	response := models.UpdatePersonResponse{
		Status: true,
	}

	// Возвращаем результат в формате XML
	fmt.Printf("Response: %+v\n", response)
	c.XML(http.StatusOK, response)
}

func (h *StorageHandler) getPersonHandler(c *gin.Context, request *models.GetPersonRequest) {
	// Получаем информацию о человеке по ID
	person, err := h.Storage.PersonRepository.GetPerson(request.ID)
	if err != nil {

		if errors.Is(err, database.ErrPersonNotFound) {
			
			return
		}

		logging.Logger.Error("Error getting person with ID", zap.Uint("ID", uint(request.ID)), zap.Error(err))
		
		return
	}

	// Если записи не найдены, формируем SOAP Fault для клиента
	if person == nil {
		fmt.Printf("No person found with ID %d\n", request.ID)

		
		return
	}

	// Если человек найден, формируем ответ
	response := models.GetPersonResponse{
		Person: *person,
	}

	// Возвращаем результат в формате XML
	fmt.Printf("Response: %+v\n", response)
	c.XML(http.StatusOK, response)
}

// Метод получения всех записей
func (h *StorageHandler) getAllPersonsHandler(c *gin.Context) {
	// Получаем все записи из базы
	persons, err := h.Storage.PersonRepository.GetAllPersons()
	if err != nil {
	
		return
	}

	response := models.GetAllPersonsResponse{
		Persons: persons,
	}

	// Возвращаем результат в формате XML
	fmt.Printf("Response: %+v\n", response)
	c.XML(http.StatusOK, response)
}

// Метод удаления записи по ID
func (h *StorageHandler) deletePersonHandler(c *gin.Context, request *models.DeletePersonRequest) {
	
	checkByID, err := h.Storage.PersonRepository.CheckPersonByID(uint(request.ID))
	if !checkByID {
	
		return
	}
	if err != nil {

		return
	}

	//Удаляем запись по ID из базы
	err = h.Storage.PersonRepository.DeletePerson(request)
	if err != nil {

		return
	}

	logging.Logger.Info("Successfully deleted person with ID", zap.Uint("ID", uint(request.ID)))
	//Формируем статус в формате SOAP

	response := models.UpdatePersonResponse{
		Status: true,
	}
	fmt.Printf("Response: %+v\n", response)
	c.XML(http.StatusOK, response)

}

// Метод поиска записей по запросу
func (h *StorageHandler) searchPersonHandler(c *gin.Context, request *models.SearchPersonRequest) {

	persons, err := h.Storage.PersonRepository.SearchPerson(request.Query)
	if err != nil {
		
		return
	}

	if len(persons) == 0 {
	
		return
	} else {
		fmt.Printf("Found persons: %+v\n", persons)
	}

	// Формируем результат в формате SOAP
	response := models.SearchPersonResponse{
		Persons: persons,
	}
	fmt.Printf("Response: %+v\n", response)
	c.XML(http.StatusOK, response)
}
