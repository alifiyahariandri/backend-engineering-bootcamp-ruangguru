package repository

import (
	"a21hc3NpZ25tZW50/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type StudentRepository interface {
	FetchAll() ([]model.Student, error)
	FetchByID(id int) (*model.Student, error)
	Store(s *model.Student) error
	Update(id int, s *model.Student) error
	Delete(id int) error
	FetchWithClass() (*[]model.StudentClass, error)
}

type studentRepoImpl struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) *studentRepoImpl {
	return &studentRepoImpl{db}
}

func (s *studentRepoImpl) FetchAll() ([]model.Student, error) {
	results := []model.Student{}
	s.db.Table("students").Select("*").Scan(&results)

	return results, nil
}

func (s *studentRepoImpl) Store(student *model.Student) error {
	result := s.db.Create(&student)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *studentRepoImpl) Update(id int, student *model.Student) error {
	var result model.Student
	s.db.Raw("SELECT * FROM students WHERE id = ?", id).Scan(&result)

	result.Name = student.Name
	result.Address = student.Address
	result.ClassId = student.ClassId
	s.db.Save(&result)

	return nil
}

func (s *studentRepoImpl) Delete(id int) error {
	s.db.Where("id = ?", id).Delete(&model.Student{})

	return nil
}

func (s *studentRepoImpl) FetchByID(id int) (*model.Student, error) {
	var result model.Student
	s.db.Raw("SELECT * FROM students WHERE id = ?", id).Scan(&result)

	if result.Name == "" {
		return &model.Student{}, errors.New("tidak ada")
	}
	return &result, nil
}

func (s *studentRepoImpl) FetchWithClass() (*[]model.StudentClass, error) {
	var result []model.StudentClass

	s.db.Table("students").Select("students.name as name, students.address as address, classes.name as class_name, classes.professor as professor, classes.room_number as room_number").Joins("join classes on students.class_id = classes.id").Scan(&result)
	fmt.Println(result)
	if result == nil {
		return &[]model.StudentClass{}, nil
	}
	return &result, nil

}
