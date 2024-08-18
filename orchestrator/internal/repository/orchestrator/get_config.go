package orchestrator

import (
	"context"
	"orchestrator/internal/config"
	"orchestrator/internal/domain"
)

var findByTypeAndService = `---
select
	order_type,
	order_service, 
	target_topic, 
	rollback_topic
from
	route_config
where
	order_type = $1 and order_service = $2
`

func (repo *OrchestratorRepositoryImpl) GetConfig(ctx context.Context) (config config.Orchestrator, err error) {
	findByTypeAndServiceStmt, err := repo.db.PrepareContext(ctx, findByTypeAndService)
	if err != nil {
		return config, err
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return config, err
	}

	orchestratorType := ctx.Value(domain.Key("type")).(string)
	orchestratorService := ctx.Value(domain.Key("service")).(string)

	row := tx.StmtContext(ctx, findByTypeAndServiceStmt).QueryRowContext(ctx, orchestratorType, orchestratorService)
	if row.Err() != nil {
		return config, err
	}

	err = row.Scan(
		&config.OrderType,
		&config.OrderService,
		&config.Topic,
		&config.RollbackTopic,
	)

	if err != nil {
		return config, err
	}

	return
}
