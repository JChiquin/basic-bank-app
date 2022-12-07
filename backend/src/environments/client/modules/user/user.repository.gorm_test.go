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
		DocumentNumber: "12345678",
		PhoneNumber:    "0123912312321",
		BirthDate:      time.Date(2020, 1, 20, 1, 2, 0, 0, time.Local),
		Email:          "some@gmail.com",
		Password:       "abc123456789",
		UserType:       constant.UserTypeClient,
	}

	connection := database.GetBankGormConnection()
	tx := connection.Begin()
	userRepo := NewUserGormRepo(tx)

	userToCreate := &entity.User{}
	copier.Copy(userToCreate, user)

	got, err := userRepo.Create(userToCreate)

	assert.NoError(t, err)
	assert.True(t, cmp.Equal(got, user, cmpopts.IgnoreFields(entity.User{}, "CreatedAt", "UpdatedAt", "ID", "Password")))
	assert.NotEqual(t, user.Password, got.Password)
	assert.NoError(t, got.CheckPassword(user.Password))
}
