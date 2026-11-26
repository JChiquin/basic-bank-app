package user

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/libs/database"
	"bank-service/src/utils/constant"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

//Setup and teardown
func TestMain(m *testing.M) {
	//setup
	database.SetupBankGormDB()
	code := m.Run() //run tests
	os.Exit(code)
}

func TestUserGormRepo(t *testing.T) {
	user := &entity.User{
		FirstName:      "Jorge",
		LastName:       "Chiquín",
		DocumentNumber: "99999999",
		PhoneNumber:    "0123912312321",
		BirthDate:      time.Date(2020, 1, 20, 1, 2, 0, 0, time.Local),
		Email:          "some@gmail.com",
		Password:       "abc123456789",
		UserType:       constant.UserTypeClient,
	}

	connection := database.GetBankGormConnection()
	tx := connection.Begin()
	defer tx.Rollback()
	userRepo := NewUserGormRepo(tx)

	userToCreate := &entity.User{}
	copier.Copy(userToCreate, user)

	got, err := userRepo.Create(userToCreate)

	assert.NoError(t, err)
	// BeforeCreate hashes the password and generates AccountNumber, so neither
	// field can be compared with the original fixture.
	assert.True(t, cmp.Equal(got, user, cmpopts.IgnoreFields(entity.User{}, "CreatedAt", "UpdatedAt", "ID", "Password", "AccountNumber")))
	assert.NotEqual(t, user.Password, got.Password)
	assert.NoError(t, got.CheckPassword(user.Password))
}
