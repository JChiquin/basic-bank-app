package contact

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/dto"

	"gorm.io/gorm"
)

/*
struct that implements IContactRepository
*/
type contactGormRepo struct {
	db *gorm.DB //Current connection
}

/*
NewContactGormRepo creates a new repo and returns IContactRepository,
so it needs to implement all its methods
*/
func NewContactGormRepo(gormDb *gorm.DB) interfaces.IContactRepository {
	return &contactGormRepo{db: gormDb}
}

func (r *contactGormRepo) IndexByUserID(contactToFilter entity.Contact, pagination *dto.Pagination) ([]entity.Contact, error) {
	return nil, nil
}
func (r *contactGormRepo) Create(newContact *entity.Contact) (*entity.Contact, error) {
	return nil, nil
}
func (r *contactGormRepo) Update(contact, updatedFields *entity.Contact) (*entity.Contact, error) {
	return nil, nil
}
func (r *contactGormRepo) Delete(contact *entity.Contact) (*entity.Contact, error) {
	return nil, nil
}
