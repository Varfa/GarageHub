package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

type SessionService struct {
	repo *repository.SessionRepository
}

func NewSessionService(
	repo *repository.SessionRepository,
) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

func (s *SessionService) Create(
	ctx context.Context,
	userID int64,
	remember bool,
) (string, time.Time, error) {
	sessionBytes := make([]byte, 32)
	_, err := rand.Read(sessionBytes)
	if err != nil {
		return "", time.Time{}, err
	}

	token := hex.EncodeToString(sessionBytes)

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	sessionDuration := 24 * time.Hour

	if remember {
		sessionDuration = 30 * 24 * time.Hour
	}

	expiresAt := time.Now().Add(
		sessionDuration,
	)
	err = s.repo.Create(ctx, userID, tokenHash, expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}
func (s *SessionService) GetUser(
	ctx context.Context,
	token string,
) (*models.User, error) {
	hash := sha256.Sum256([]byte(token))

	tokenHash := hex.EncodeToString(hash[:])
	user, err := s.repo.GetUserByTokenHash(
		ctx,
		tokenHash,
	)
	return user, err
}
func (s *SessionService) Delete(
	ctx context.Context,
	token string,
) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	err := s.repo.Delete(ctx, tokenHash)
	return err
}
