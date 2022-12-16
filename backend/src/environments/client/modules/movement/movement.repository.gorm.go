package movement

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"

	"gorm.io/gorm"
)

/*
struct that implements IMovementRepository
*/
type movementGormRepo struct {
	db *gorm.DB //Current connection
}

/*
NewMovementGormRepo creates a new repo and returns IMovementRepository,
so it needs to implement all its methods
*/
func NewMovementGormRepo(gormDb *gorm.DB) interfaces.IMovementRepository {
	return &movementGormRepo{db: gormDb}
}

func (r *movementGormRepo) IndexByUserID(userID int) ([]entity.Movement, error) {
	movements := []entity.Movement{}
	err := r.db.Model(entity.Movement{}).
		Where(entity.Movement{UserID: userID}).
		Order("created_at desc").
		Find(&movements).Error
	if err != nil {
		return nil, err
	}
	return movements, nil
}
