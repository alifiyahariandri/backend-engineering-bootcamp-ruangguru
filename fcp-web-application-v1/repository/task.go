package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Store(task *model.Task) error
	Update(id int, task *model.Task) error
	Delete(id int) error
	GetByID(id int) (*model.Task, error)
	GetList() ([]model.Task, error)
	GetTaskCategory(id int) ([]model.TaskCategory, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *taskRepository {
	return &taskRepository{db}
}

func (t *taskRepository) Store(task *model.Task) error {
	err := t.db.Create(task).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *taskRepository) Update(id int, task *model.Task) error {
	err := t.db.Model(&task).Where("id = ?", id).Update("title", task.Title).Error
	if err != nil {
		return err
	}
	err = t.db.Model(&task).Where("id = ?", id).Update("deadline", task.Deadline).Error
	if err != nil {
		return err
	}
	err = t.db.Model(&task).Where("id = ?", id).Update("priority", task.Priority).Error
	if err != nil {
		return err
	}
	err = t.db.Model(&task).Where("id = ?", id).Update("category_id", task.CategoryID).Error
	if err != nil {
		return err
	}
	err = t.db.Model(&task).Where("id = ?", id).Update("status", task.Status).Error
	if err != nil {
		return err
	}

	return nil

}

func (t *taskRepository) Delete(id int) error {
	task, err := t.GetByID(id)

	err = t.db.Where("id = ?", id).Delete(&task).Error

	if err != nil {
		return err
	}

	return nil
}

func (t *taskRepository) GetByID(id int) (*model.Task, error) {
	var task model.Task
	err := t.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (t *taskRepository) GetList() ([]model.Task, error) {
	var tasks []model.Task
	t.db.Find(&tasks)

	return tasks, nil
}

func (t *taskRepository) GetTaskCategory(id int) ([]model.TaskCategory, error) {
	var tasks []model.TaskCategory
	t.db.Table("tasks").Select("tasks.ID, tasks.title, categories.name as category").Joins("join categories on tasks.category_id = categories.id").Where("tasks.id = $1", id).Scan(&tasks)
	// t.db.Table("tasks").Select("tasks.ID, tasks.title, categories.name as category").Joins("join categories on tasks.category_id = categories.id").Where("categories.id = $1", id).Scan(&tasks)

	return tasks, nil
}
