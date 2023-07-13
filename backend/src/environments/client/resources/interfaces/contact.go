package interfaces

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/libs/dto"
	"net/http"
)

/*
IContactController methods to handle http requests
*/
type IContactController interface {
	Index(response http.ResponseWriter, request *http.Request)
	Create(response http.ResponseWriter, request *http.Request)
	Update(response http.ResponseWriter, request *http.Request)
	Delete(response http.ResponseWriter, request *http.Request)
	GetOne(response http.ResponseWriter, request *http.Request)
}

/*
IContactService methods to implement the bussiness logic
*/
type IContactService interface {
	IndexByUserID(filterContacts *dto.FilterContacts, pagination *dto.Pagination) ([]entity.Contact, error)
	Create(createContact *dto.CreateContact) (*entity.Contact, error)
	Update(updateContact *dto.UpdateContact) (*entity.Contact, error)
	Delete(filterOneContact *dto.FilterOneContact) (*entity.Contact, error)
	GetOne(filterOneContact *dto.FilterOneContact) (*entity.Contact, error)
}

/*
IContactRepository methods to interact with contact entity, independent of ORM
*/
type IContactRepository interface {
	IndexByUserID(filterContacts *dto.FilterContacts, pagination *dto.Pagination) ([]entity.Contact, error)
	Create(newContact *entity.Contact) (*entity.Contact, error)
	Update(contact, updatedFields *entity.Contact) (*entity.Contact, error)
	Delete(contact *entity.Contact) (*entity.Contact, error)
	FindByAlias(alias string, userID int) (*entity.Contact, error)
	FindByAccountNumber(accountNumber string, userID int) (*entity.Contact, error)
	FindByID(id, userID int) (*entity.Contact, error)
	FindFullContactByID(id, userID int) (*entity.Contact, error)
}
