create table if not exists users
(
    id serial primary key,
    login varchar(255) not null,
    password varchar(255) not null,
    created_at timestamp default current_timestamp
    );

create index if not exists idx_users_login on users(login);
