package data

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db}
}

func (r *SessionRepository) GetByToken(token string) (*Session, error) {
	var session Session
	result := r.db.Where("token = ?", token).First(&session)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (r *SessionRepository) Create(req *SessionCreate) (*Session, error) {
	guid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		Token:     guid.String(),
		UserAgent: req.UserAgent,
		IP:        req.IP,
	}

	result := r.db.Create(session)
	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}
