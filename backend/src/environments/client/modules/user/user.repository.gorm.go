package user

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	myErrors "bank-service/src/libs/errors"
	"errors"

	"gorm.io/gorm"
)

//struct that implements IUserRepository

type userGormRepo struct {
	db *gorm.DB //Current connection
}

/*
NewUserGormRepo creates a new repo and returns IUserRepository,
so needs to implement all its methods
*/
func NewUserGormRepo(gormDb *gorm.DB) interfaces.IUserRepository {
	rUser := userGormRepo{gormDb}
	return &rUser
}

/*
Create receives an user entity, sends to database and returns it
If there is an error, returns it as a second result
*/
func (r *userGormRepo) Create(user *entity.User) (*entity.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

/*
FindByDocumentNumber receives the document number and tries to find one user with that value
If there is an error, returns it as a second result
*/
func (r *userGormRepo) FindByDocumentNumber(documentNumber string) (*entity.User, error) {
	return r.findByAttributes(entity.User{DocumentNumber: documentNumber})
}

/*
FindByEmail receives the email and tries to find one user with that email
If there is an error, returns it as a second result
*/
func (r *userGormRepo) FindByEmail(email string) (*entity.User, error) {
	return r.findByAttributes(entity.User{Email: email})
}

/*
FindByID receives the user ID and tries to find one user with that value
If there is an error, returns it as a second result
*/
func (r *userGormRepo) FindByID(userID int) (*entity.User, error) {
	return r.findByAttributes(entity.User{ID: userID})
}

/*
FindByID receives the user account number and tries to find one user with that value
If there is an error, returns it as a second result
*/
func (r *userGormRepo) FindByAccountNumber(accountNumber string) (*entity.User, error) {
	return r.findByAttributes(entity.User{AccountNumber: accountNumber})
}

/*
GetBalance receives the user id and find its last movement to get the current balance
If there is an error, returns it as a second result
*/
func (r *userGormRepo) GetBalance(userID int) (float64, error) {
	lastMovement := &entity.Movement{}
	err := r.db.
		Select("balance").
		Where(entity.Movement{UserID: userID}).
		Order("created_at DESC").
		Take(lastMovement).
		Error

	//If no transactions were found, the balance will be initial value (0)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return lastMovement.Balance, nil
}

func (r *userGormRepo) findByAttributes(userFilter entity.User) (*entity.User, error) {
	user := &entity.User{}

	err := r.db.Take(user, userFilter).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, myErrors.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userGormRepo) UpdatePassword(userID int, newPlainPassword string) error {
	result := r.db.
		Model(&entity.User{ID: userID}).
		Updates(entity.User{Password: newPlainPassword})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return myErrors.ErrNotFound
	}

	return nil
}
