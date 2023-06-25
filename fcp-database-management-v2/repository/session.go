package repository

import (
	"a21hc3NpZ25tZW50/model"
	"errors"

	"gorm.io/gorm"
)

type SessionsRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailName(name string) error
	SessionAvailToken(token string) (model.Session, error)
}

type sessionsRepoImpl struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *sessionsRepoImpl {
	return &sessionsRepoImpl{db}
}

func (s *sessionsRepoImpl) AddSessions(session model.Session) error {
	result := s.db.Create(&session)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *sessionsRepoImpl) DeleteSession(token string) error {
	s.db.Where("token = ?", token).Delete(&model.Session{})

	return nil
}

func (s *sessionsRepoImpl) UpdateSessions(session model.Session) error {
	var result model.Session
	s.db.Raw("SELECT * FROM sessions WHERE username = ?", session.Username).Scan(&result)

	result.Token = session.Token
	result.Expiry = session.Expiry
	s.db.Save(&result)

	return nil
}

func (s *sessionsRepoImpl) SessionAvailName(name string) error {
	var result model.Session
	s.db.Raw("SELECT * FROM sessions WHERE username = ?", name).Scan(&result)

	if result.Username == "" {
		return errors.New("tidak ada")
	}
	return nil
}

func (s *sessionsRepoImpl) SessionAvailToken(token string) (model.Session, error) {
	var result model.Session
	s.db.Raw("SELECT * FROM sessions WHERE token = ?", token).Scan(&result)

	if result.Username == "" {
		return model.Session{}, errors.New("tidak ada")
	}
	return result, nil
}
