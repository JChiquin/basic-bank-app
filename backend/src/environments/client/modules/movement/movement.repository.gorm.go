package movement

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/dto"

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

func (r *movementGormRepo) IndexByUserID(movementToFilter entity.Movement, pagination *dto.Pagination) ([]entity.Movement, error) {
	movements := []entity.Movement{}
	err := r.db.Model(entity.Movement{}).
		Where(movementToFilter).
		Count(&pagination.TotalCount).
		Offset(pagination.Offset()).
		Limit(pagination.PageSize).
		Order("created_at desc").
		Find(&movements).Error
	if err != nil {
		return nil, err
	}
	return movements, nil
}
