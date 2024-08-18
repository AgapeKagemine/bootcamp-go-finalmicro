package orchestrator

import (
	"context"
	"encoding/json"
	"orchestrator/internal/domain"
	"orchestrator/internal/domain/orchestrator"
)

var getByIdFailed = `---
with max_retries as (
    select
        transaction_id,
        max(retries) as max_retries
    from
        transaction
    where
        response_code > 201
    group by
        transaction_id
)
select
    t.transaction_id,
    t.transaction_datetime,
    t.order_type,
    t.order_service,
    t.retries,
    t.response_code,
    t.response_message,
    t.request_body 
from
    transaction t
join
    max_retries mr
	on
    t.transaction_id = mr.transaction_id
	and
    t.retries = mr.max_retries
where
    t.response_code > 201
	and
	t.transaction_id = $1;
`

func (repo *OrchestratorRepositoryImpl) GetByIdFailed(ctx context.Context) (trx orchestrator.Message, err error) {
	getByIdFailedStmt, err := repo.db.PrepareContext(ctx, getByIdFailed)
	if err != nil {
		return trx, err
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return trx, err
	}

	trx_req := ctx.Value(domain.Key("request")).(orchestrator.Message)

	row := tx.StmtContext(ctx, getByIdFailedStmt).QueryRowContext(ctx, trx_req.Header.TransactionId)
	err = row.Scan(
		&trx.Header.TransactionId,
		&trx.Header.TransactionDateTime,
		&trx.Header.OrderType,
		&trx.Header.OrderService,
		&trx.Header.Retries,
		&trx.Header.ResponseCode,
		&trx.Header.ResponseMessage,
		&trx.Body,
	)

	if err != nil {
		return trx, err
	}

	body := make(map[string]interface{}, 0)
	err = json.Unmarshal(trx.Body.([]byte), &body)
	if err != nil {
		return trx, err
	}
	trx.Body = body

	return
}
