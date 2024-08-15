drop table if exists route_config;

create table route_config (
    order_type VARCHAR(255) not null,
    order_service VARCHAR(255) not null,
    target_topic VARCHAR(255) not null
);

insert into 
	route_config (order_type, order_service, target_topic)
values
	('ORDER BOOK', 'START', 'user-validate'),
	('ORDER BOOK', 'user-validate', 'book-validate'),
	('ORDER BOOK', 'book-validate', 'payment-deduct-balance'),
	('ORDER BOOK', 'payment-deduct-balance', 'FINISH');
	
insert into 
	route_config (order_type, order_service, target_topic)
values
	('Transfer Balance', 'START', 'validateUser'),
	('Transfer Balance', 'validateUser', 'deductBalance'),
	('Transfer Balance', 'deductBalance', 'addBalance'),
	('Transfer Balance', 'addBalance', 'FINISH');

select
	*
from
	route_config
where
	order_type = 'ORDER BOOK' and order_service = 'payment'
limit 
	1;