package orchestrator

import (
	"context"
	"encoding/json"
	"orchestrator/internal/domain/orchestrator"
)

var getAllFailed = `---
select
	transaction_id,
	transaction_datetime,
	order_type,
	order_service,
	retries,
	response_code,
	response_message,
	request_body
from
	transaction
where
	response_code > 201
`

func (repo *OrchestratorRepositoryImpl) GetAllFailed(ctx context.Context) (trx []orchestrator.Message, err error) {
	getAllFailedStmt, err := repo.db.PrepareContext(ctx, getAllFailed)
	if err != nil {
		return nil, err
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.StmtContext(ctx, getAllFailedStmt).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		msg := &orchestrator.Message{}
		err := rows.Scan(
			&msg.Header.TransactionId,
			&msg.Header.TransactionDateTime,
			&msg.Header.OrderType,
			&msg.Header.OrderService,
			&msg.Header.Retries,
			&msg.Header.ResponseCode,
			&msg.Header.ResponseMessage,
			&msg.Body,
		)
		if err != nil {
			return nil, err
		}

		body := make(map[string]interface{}, 0)
		err = json.Unmarshal(msg.Body.([]byte), &body)
		if err != nil {
			return nil, err
		}
		msg.Body = body

		trx = append(trx, *msg)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return
}
