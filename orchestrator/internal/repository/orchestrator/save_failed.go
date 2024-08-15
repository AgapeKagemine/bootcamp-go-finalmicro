package orchestrator

import (
	"context"
	"fmt"
	"orchestrator/internal/domain"
	"orchestrator/internal/domain/orchestrator"
)

var saveFailed = `---
insert into 
	failed_transaction (transaction_id, transaction_datetime, order_type, order_service, retries, response_code, response_message, request_body)
values
	($1, $2, $3, $4, $5, $6, $7, $8)
`

func (repo *OrchestratorRepositoryImpl) SaveFailed(ctx context.Context) error {
	request := ctx.Value(domain.Key("request")).(orchestrator.Message)

	saveFailedStmt, err := repo.db.PrepareContext(ctx, saveFailed)
	if err != nil {
		return err
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	row := tx.StmtContext(ctx, saveFailedStmt).QueryRowContext(ctx, request.Header.TransactionId, request.Header.TransactionDateTime, request.Header.OrderType, request.Header.OrderService, request.Header.Retries, request.Header.ResponseCode, request.Header.ResponseMessage, request.Body)
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
