package repository

import (
	"auth-service/internal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByID(id int64) (*entity.User, error)
	FindByIDs(ids []int64) ([]entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id int64) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByIDs(ids []int64) ([]entity.User, error) {
	var users []entity.User
	// GORM akan otomatis mengonversi slice ids menjadi query:
	// SELECT * FROM users WHERE id IN (1, 2, 3)
	err := r.db.Where("id IN ?", ids).Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}
