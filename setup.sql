create table rotations (
    id serial primary key,
    banner_id bigint not null,
    slot_id bigint not null,
    description text not null,
    create_at timestamp not null
);
create index slot_idx on rotations (slot_id);
create index banner_idx on rotations (banner_id);
create table statistics (
    id serial primary key,
    type bigint not null,
    banner_id bigint not null,
    slot_id bigint not null,
    group_id bigint not null,
    create_at timestamp not null
);
create index banner_idx on statistics (banner_id);
create index slot_idx on statistics (slot_id);
create index group_idx on statistics (group_id);
