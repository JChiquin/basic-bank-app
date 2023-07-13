package contact

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/dto"
	"bank-service/src/libs/errors"

	goerrors "errors"

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

func (r *contactGormRepo) IndexByUserID(filter *dto.FilterContacts, pagination *dto.Pagination) ([]entity.Contact, error) {
	contacts := []entity.Contact{}
	err := r.db.Model(contacts).
		Where(entity.Contact{UserID: filter.UserID}).
		Where(`alias ILIKE ?`, filter.LikeValue()).
		Count(&pagination.TotalCount).
		Offset(pagination.Offset()).
		Limit(pagination.PageSize).
		Order("alias ASC").
		Find(&contacts).Error
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *contactGormRepo) Create(newContact *entity.Contact) (*entity.Contact, error) {
	err := r.db.Model(entity.Contact{}).
		Create(&newContact).Error
	if err != nil {
		return nil, err
	}
	return newContact, nil
}

func (r *contactGormRepo) Update(contact, updatedFields *entity.Contact) (*entity.Contact, error) {
	result := r.db.Model(&contact).
		Updates(&updatedFields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.ErrNotFound
	}
	return contact, nil
}

func (r *contactGormRepo) Delete(contact *entity.Contact) (*entity.Contact, error) {
	result := r.db.Delete(contact)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.ErrNotFound
	}
	return contact, nil
}

func (r *contactGormRepo) FindByAlias(alias string, userID int) (*entity.Contact, error) {
	return r.findByAttributes(entity.Contact{Alias: alias, UserID: userID}, false)
}

func (r *contactGormRepo) FindByAccountNumber(accountNumber string, userID int) (*entity.Contact, error) {
	return r.findByAttributes(entity.Contact{AccountNumber: accountNumber, UserID: userID}, false)
}

func (r *contactGormRepo) FindByID(id, userID int) (*entity.Contact, error) {
	return r.findByAttributes(entity.Contact{ID: id, UserID: userID}, false)
}

func (r *contactGormRepo) FindFullContactByID(id, userID int) (*entity.Contact, error) {
	return r.findByAttributes(entity.Contact{ID: id, UserID: userID}, true)
}

func (r *contactGormRepo) findByAttributes(contactFilter entity.Contact, preloadUser bool) (*entity.Contact, error) {
	contact := &entity.Contact{}
	db := r.db.Model(entity.Contact{})

	if preloadUser {
		db = db.Preload("User")
	}

	err := db.Take(contact, contactFilter).Error
	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return contact, nil
}
