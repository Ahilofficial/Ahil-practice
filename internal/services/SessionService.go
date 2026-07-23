package services

import (
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
)

type SessionService struct {
	sessionrepo *repository.SessionRepository
}

func NewSessionService(sessionrepo *repository.SessionRepository) *SessionService {
	return &SessionService{sessionrepo: sessionrepo}
}
func (s *SessionService) CreateSession(userID uint, platform string, session_id string, access_token string, refresh_token string) (*model.Session, error) {
	session := &model.Session{
		UserID:   userID,
		IsActive: true,
		Platform: platform,
		SessionID: session_id,
		AccessToken: access_token,
		RefreshToken: refresh_token,
	}

	err := s.sessionrepo.CreateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
