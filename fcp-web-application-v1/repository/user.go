package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	// fmt.Println("AAAAAAAAA")
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		// fmt.Println(err)
		return model.User{}, nil
	}

	return user, nil
}

func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) GetUserTaskCategory() ([]model.UserTaskCategory, error) {

	var usertaskcat []model.UserTaskCategory
	r.db.Table("users").Select("users.id, users.fullname, users.email, tasks.title as task, tasks.deadline, tasks.priority, tasks.status, categories.name as category").Joins("join tasks on users.id = tasks.user_id join categories on tasks.category_id = categories.id").Scan(&usertaskcat)

	return usertaskcat, nil
}
