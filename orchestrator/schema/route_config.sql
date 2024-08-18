-- Active: 1717482674968@@127.0.0.1@5432@gomicro-orchestrator
drop table if exists route_config;
drop table if exists transaction;
drop table if exists failed_transaction;

create table route_config (
    order_type VARCHAR(255) not null,
    order_service VARCHAR(255) not null,
    target_topic VARCHAR(255) not null,
    rollback_topic VARCHAR(255) not null,
    created_at timestamptz default timezone('Asia/Jakarta', current_timestamp)
);

insert into 
	route_config (order_type, order_service, target_topic, rollback_topic)
values
	('ORDER BOOK', 'START', 'user-validate', 'FINISH'),
	('ORDER BOOK', 'user-validate', 'book-validate', 'FINISH'),
	('ORDER BOOK', 'book-validate', 'payment-deduct-balance', 'FINISH'),
	('ORDER BOOK', 'payment-deduct-balance', 'FINISH', 'ROLLBACK'),
	('ORDER BOOK', 'ROLLBACK', 'rollback-book', 'FINISH'),
	('ORDER BOOK', 'rollback-book', 'FINISH', 'FINISH');

create table transaction (
	transaction_id VARCHAR(255) not null,
	transaction_datetime timestamptz not null,
	order_type VARCHAR(255) not null,
	order_service VARCHAR(255) not null,
	retries integer not null,
	response_code integer not null,
	response_message VARCHAR(255) not null,
	request_body JSONB not null,
	created_at timestamptz default timezone('Asia/Jakarta', current_timestamp),
	updated_at timestamptz default timezone('Asia/Jakarta', current_timestamp)
);

select
	target_topic
from
	route_config
where
	order_type = 'ORDER BOOK'
	and 
	order_service = 'ROLLBACK';
	
select
	*
from
	transaction;

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
	t.transaction_id = '019160eb-9efd-7684-b71c-eea5f283bc7c';
	
SELECT
	*
from
	transaction
where
	response_code > 201
	and
	transaction_id = '019155c2-8cae-7cb4-8ac7-68a0dbcec24d';