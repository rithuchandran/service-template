create table regions (
  id bigint unique not null ,
  name text,
  data jsonb
);