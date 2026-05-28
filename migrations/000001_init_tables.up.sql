create table if not exists users
(
    id serial primary key,
    login varchar(255) not null,
    password varchar(255) not null,
    created_at timestamp default current_timestamp
    );

create index if not exists idx_users_login on users(login);

create table if not exists orders
(
    id serial primary key,
    number varchar(255) not null,
    status varchar(255) not null,
    uploaded_at timestamp default current_timestamp,
    accrual numeric(1000, 2),
    user_id integer not null
    );

create unique index if not exists idx_orders_number on orders(number);

create table if not exists balances
(
    id serial primary key,
    current numeric(1000, 2) not null,
    withdrawn numeric(1000, 2) not null,
    user_id integer not null
    );

create unique index if not exists idx_balances_user_id on balances(user_id);

create table if not exists withdrawns
(
    id serial primary key,
    sum numeric(1000, 2) not null,
    order_num varchar(255) not null,
    user_id integer not null,
    processed_at timestamp default current_timestamp
    );

create index if not exists idx_withdrawns_user_id on withdrawns(user_id);