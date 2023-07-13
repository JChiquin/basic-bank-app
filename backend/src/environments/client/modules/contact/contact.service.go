package contact

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/dto"
)

/*
struct that implements IContactService
*/
type contactService struct {
	rContact interfaces.IContactRepository
	rUser    interfaces.IUserRepository
}

/*
NewContactService creates a new service, receives repository by dependency injection
and returns IContactService, so it needs to implement all its methods
*/
func NewContactService(rContact interfaces.IContactRepository, rUser interfaces.IUserRepository) interfaces.IContactService {
	return &contactService{rContact, rUser}
}

func (s *contactService) IndexByUserID(filterContacts *dto.FilterContacts, pagination *dto.Pagination) ([]entity.Contact, error) {
	return nil, nil
}
func (s *contactService) Create(createContact *dto.CreateContact) (*entity.Contact, error) {
	return nil, nil
}
func (s *contactService) Update(updateContact *dto.UpdateContact) (*entity.Contact, error) {
	return nil, nil
}
func (s *contactService) Delete(contactID, userID int) (*entity.Contact, error) {
	return nil, nil
}
