-- +goose Up
create table orders (
    id serial primary key,
    order_uuid text not null unique,
    user_uuid text not null,
    total_price numeric not null,
    transaction_uuid text,
    payment_method text,
    status text not null
);

create table order_parts (
    part_uuid text not null,
    order_uuid text not null,
    primary key (order_uuid, part_uuid)
);

create index idx_order_parts_part_uuid on order_parts (part_uuid);

-- +goose Down
drop table if exists order_parts;
drop table if exists orders;

