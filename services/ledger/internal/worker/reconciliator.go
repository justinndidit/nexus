package worker

import (
	"github.com/justinndidit/nexus/ledger/internal/platform/broker"
	"github.com/rs/zerolog"
)

type Reconciliator struct {
	logger    *zerolog.Logger
	publisher broker.Publisher
}

func NewReconciliator(logger *zerolog.Logger, pub broker.Publisher) *Reconciliator {
	return &Reconciliator{
		logger:    logger,
		publisher: pub,
	}
}

func (r *Reconciliator) IntegrityChecker() error {
	return nil
}

func (r *Reconciliator) Reconcile() error {
	return nil
}
