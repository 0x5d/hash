create table if not exists url_ids(
  id serial primary key,
  url_id bigint not null check (url_id >= 0)
);

insert into url_ids(id, url_id) values (default, 0);

create table if not exists urls(
  id serial primary key,
  url varchar(2048) not null,
  short_key varchar(2048) not null,
  enabled boolean not null default true
);
