package orchestrator

import (
	"context"
	"database/sql"
	"orchestrator/internal/config"
	"orchestrator/internal/domain/orchestrator"
)

type OrchestratorRepository[C context.Context, T orchestrator.Message, E error] interface {
	GetConfig(C) (config.OrchestratorConfig, E)
	Save(C) E
	GetAll(C) ([]T, E)
	GetAllFailed(C) ([]T, E)
	GetByIdFailed(C) (T, E)
}

type OrchestratorRepositoryImpl struct {
	db *sql.DB
}

func NewOrchestratorRepository(db *sql.DB) OrchestratorRepository[context.Context, orchestrator.Message, error] {
	return &OrchestratorRepositoryImpl{
		db: db,
	}
}
