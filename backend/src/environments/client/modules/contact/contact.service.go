package contact

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/dto"
	"bank-service/src/libs/errors"
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
	if err := filterContacts.Validate(); err != nil {
		return nil, err
	}

	listContact, err := s.rContact.IndexByUserID(filterContacts, pagination)
	if err != nil {
		return nil, err
	}

	return listContact, nil
}
func (s *contactService) Create(createContact *dto.CreateContact) (*entity.Contact, error) {
	if err := createContact.Validate(); err != nil {
		return nil, err
	}

	user, err := s.rUser.FindByAccountNumber(createContact.AccountNumber)
	if err != nil {
		return nil, err
	}

	if user.ID == createContact.UserID {
		return nil, errors.ErrSameContact
	}

	contact, _ := s.rContact.FindByAccountNumber(createContact.AccountNumber, createContact.UserID)
	if contact != nil {
		return nil, errors.ErrDuplicatedContact
	}
	contact, _ = s.rContact.FindByAlias(createContact.Alias, createContact.UserID)
	if contact != nil {
		return nil, errors.ErrAliasInUse
	}

	newContact, err := s.rContact.Create(createContact.ParseToEntity())
	if err != nil {
		return nil, err
	}

	return newContact, nil
}

func (s *contactService) Update(updateContact *dto.UpdateContact) (*entity.Contact, error) {
	if err := updateContact.Validate(); err != nil {
		return nil, err
	}

	contactToUpdate, err := s.rContact.FindByID(updateContact.ContactID, updateContact.UserID)
	if err != nil {
		return nil, err
	}

	contact, _ := s.rContact.FindByAlias(updateContact.Alias, updateContact.UserID)
	if contact != nil && contact.ID != updateContact.ContactID {
		return nil, errors.ErrAliasInUse
	}

	updatedContact, err := s.rContact.Update(contactToUpdate, updateContact.ParseToEntity())
	if err != nil {
		return nil, err
	}

	return updatedContact, nil
}

func (s *contactService) Delete(filterOneContact *dto.FilterOneContact) (*entity.Contact, error) {
	if err := filterOneContact.Validate(); err != nil {
		return nil, err
	}

	contactToDelete, err := s.rContact.FindByID(filterOneContact.ContactID, filterOneContact.UserID)
	if err != nil {
		return nil, err
	}

	deletedContact, err := s.rContact.Delete(contactToDelete)
	if err != nil {
		return nil, err
	}

	return deletedContact, nil
}

func (s *contactService) GetOne(filterOneContact *dto.FilterOneContact) (*entity.Contact, error) {
	if err := filterOneContact.Validate(); err != nil {
		return nil, err
	}

	contact, err := s.rContact.FindFullContactByID(filterOneContact.ContactID, filterOneContact.UserID)
	if err != nil {
		return nil, err
	}

	return contact, nil
}
