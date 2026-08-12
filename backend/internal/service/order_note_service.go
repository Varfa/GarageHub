package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Varfa/GarageHub/internal/models"
	"github.com/Varfa/GarageHub/internal/repository"
)

var ErrOrderNoteTextRequired = errors.New(
	"order_note.text_required",
)

type OrderNoteService struct {
	repo *repository.OrderNoteRepository
}

func NewOrderNoteService(
	repo *repository.OrderNoteRepository,
) *OrderNoteService {
	return &OrderNoteService{
		repo: repo,
	}
}

func (s *OrderNoteService) Create(
	ctx context.Context,
	note models.OrderNote,
) error {
	note.Text = strings.TrimSpace(
		note.Text,
	)
	if note.OrderID <= 0 {
		return ErrOrderInvalidID
	}
	if note.Text == "" {
		return ErrOrderNoteTextRequired
	}

	return s.repo.Create(
		ctx,
		note,
	)
}

func (s *OrderNoteService) ListByOrderID(
	ctx context.Context,
	orderID int64,
) ([]models.OrderNote, error) {

	if orderID <= 0 {
		return nil, ErrOrderInvalidID
	}

	return s.repo.ListByOrderID(
		ctx, orderID,
	)

}
