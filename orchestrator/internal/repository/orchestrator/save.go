package orchestrator

import (
	"context"
	"fmt"
	"orchestrator/internal/domain"
	"orchestrator/internal/domain/orchestrator"
)

var save = `---
insert into 
	transaction (transaction_id, transaction_datetime, order_type, order_service, retries, response_code, response_message, request_body)
values
	($1, $2, $3, $4, $5, $6, $7, $8)
`

func (repo *OrchestratorRepositoryImpl) Save(ctx context.Context) error {
	request := ctx.Value(domain.Key("request")).(orchestrator.Message)

	saveStmt, err := repo.db.PrepareContext(ctx, save)
	if err != nil {
		return err
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	row := tx.StmtContext(ctx, saveStmt).QueryRowContext(ctx, request.Header.TransactionId, request.Header.TransactionDateTime, request.Header.OrderType, request.Header.OrderService, request.Header.Retries, request.Header.ResponseCode, request.Header.ResponseMessage, request.Body)
	if row.Err() != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save failed_transaction: %w", row.Err())
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
