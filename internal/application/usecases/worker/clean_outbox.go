package workeruc

import (
	"context"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/interfaces"
)

type CleanOutboxUsecase struct {
	cleaner interfaces.OutboxCleaner
}

func NewCleanOutboxUsecasee(cleaner interfaces.OutboxCleaner) *CleanOutboxUsecase {
	return &CleanOutboxUsecase{
		cleaner: cleaner,
	}
}

func (uc CleanOutboxUsecase) Execute(ctx context.Context) error {
	return uc.cleaner.Clean(ctx)
}
