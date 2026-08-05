package worker

import (
	"github.com/justinndidit/nexus/ledger/internal/ledger"
	"github.com/justinndidit/nexus/ledger/internal/platform/broker"
	"github.com/rs/zerolog"
)

type Reconciliator struct {
	logger    *zerolog.Logger
	repo      ledger.Repository
	publisher broker.Publisher
}

func NewReconciliator(logger *zerolog.Logger, repo ledger.Repository, pub broker.Publisher) *Reconciliator {
	return &Reconciliator{
		logger:    logger,
		repo:      repo,
		publisher: pub,
	}
}

func (r *Reconciliator) IntegrityChecker() error {
	return nil
}

func (r *Reconciliator) Reconcile() error {
	return nil
}
