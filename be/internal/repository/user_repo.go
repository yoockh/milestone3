package repository

import (
	"context"
	"milestone3/be/internal/entity"

	"gorm.io/gorm"
)

type UserRepo struct {
	ctx context.Context
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB, ctx context.Context) *UserRepo {
	return &UserRepo{db: db, ctx: ctx}
}

func (ur *UserRepo) Create(user *entity.Users) error {
	if err := ur.db.WithContext(ur.ctx).Omit("RoleId").Create(user).Error; err != nil {
		return err
	}
	
	return nil
}

func (ur *UserRepo) GetByEmail(email string) (user entity.Users, err error) {
	if err := ur.db.WithContext(ur.ctx).Preload("Role").First(&user, "email = ?", email).Error; err != nil {
		return entity.Users{}, err
	}

	return user, nil
}
func (ur *UserRepo) GetById(id int) (user entity.Users, err error) {
	if err := ur.db.WithContext(ur.ctx).First(&user, "id = ?", id).Error; err != nil {
		return entity.Users{}, err
	}

	return user, nil
}
// reset password?
// validation?