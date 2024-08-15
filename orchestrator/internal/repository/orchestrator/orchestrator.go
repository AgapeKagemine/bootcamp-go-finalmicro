package orchestrator

import (
	"context"
	"database/sql"
	"orchestrator/internal/config"
)

type OrchestratorRepository[C context.Context] interface {
	GetConfig(C) (config.Orchestrator, error)
	Save(C) error
	SaveFailed(C) error
}

type OrchestratorRepositoryImpl struct {
	db *sql.DB
}

func NewOrchestratorRepository(db *sql.DB) OrchestratorRepository[context.Context] {
	return &OrchestratorRepositoryImpl{
		db: db,
	}
}
