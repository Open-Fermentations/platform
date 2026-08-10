begin;

  create extension if not exists "uuid-ossp";

  create table "user"
  (
    id uuid primary key default uuid_generate_v4(),
    username varchar not null,
    -- TODO: if a password is always hashed then it should be the same length all the time?
    password varchar not null, 
    active boolean default true not null,
    created timestamptz default current_timestamp not null,
    modified timestamptz default current_timestamp not null
  );

  create unique index unique_username_constraint on "user"(username);

  create table "role"
  (
    id uuid primary key default uuid_generate_v4(),
    "name" varchar(20) not null
  );

  create unique index unique_role_name_constraint on "role"("name");

  create table "permission"
  (
    id uuid primary key default uuid_generate_v4(),
    "name" varchar(20) not null
  );

  create unique index unique_permission_name_constraint on "permission"("name");

  create table "role_permission"
  (
    id uuid primary key default uuid_generate_v4(),
    role_id uuid not null,
    permission_id uuid not null,
    constraint fk_role_permission_role
      foreign key(role_id)
      references "role"(id)
      on delete cascade,
    constraint fk_role_permission_permission
      foreign key(permission_id)
      references "permission"(id)
      on delete cascade
  );

  -- TODO: add constraint for unique role_id and permission_id

  create table "user_role"
  (
    id uuid primary key default uuid_generate_v4(),
    user_id uuid not null,
    role_id uuid not null,
    constraint fk_user_role_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade,
    constraint fk_user_role_role
      foreign key(role_id)
      references "role"(id)
      on delete cascade
  );

  -- TODO: add constraint for unique user_id and role_id

  create table "user_permission"
  (
    id uuid primary key default uuid_generate_v4(),
    user_id uuid not null,
    permission_id uuid not null,
    constraint fk_user_permission_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade,
    constraint fk_user_permission_permission
      foreign key(permission_id)
      references "permission"(id)
      on delete cascade
  );

  -- TODO: create constraint for unique user_id and permission_id

  create table "batch"
  (
    id uuid primary key default uuid_generate_v4(),
    "name" varchar not null,
    user_id uuid not null,
    created timestamptz default current_timestamp not null,
    modified timestamptz default current_timestamp not null,
    constraint fk_batch_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade
  );

  create type reading_type_enum as enum
  ('temperature', 'gravity', 'volume', 'pressure', 'ph');

  create table "device"
  (
    id uuid primary key default uuid_generate_v4(),
    "name" varchar not null,
    mac_address macaddr8 not null,
    user_id uuid not null,
    created timestamptz default current_timestamp not null,
    modified timestamptz default current_timestamp not null,
    constraint fk_device_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade
  );

  create unique index unique_mac_address_constraint on "device"(mac_address);

  create table "batch_device"
  (
    id uuid primary key default uuid_generate_v4(),
    batch_id uuid not null,
    device_id uuid not null,
    created timestamptz default current_timestamp not null,
    modified timestamptz default current_timestamp not null,
    constraint fk_batch_device_batch
      foreign key(batch_id)
      references "batch"(id)
      on delete cascade,
    constraint fk_batch_device_device
      foreign key(device_id)
      references "device"(id)
      on delete cascade
  );

  create table "device_capability"
  (
    id uuid primary key default uuid_generate_v4(),
    device_id uuid not null,
    capability reading_type_enum not null,
    created timestamptz default current_timestamp not null,
    modified timestamptz default current_timestamp not null,
    constraint fk_device_capability_device
      foreign key(device_id)
      references "device"(id)
      on delete cascade
  );

  create table "reading"
  (
    id uuid primary key default uuid_generate_v4(),
    reading_type reading_type_enum not null,
    "value" double precision not null,
    device_id uuid not null,
    user_id uuid not null,
    created timestamptz default current_timestamp not null,
    modified timestamptz default current_timestamp not null,
    constraint fk_reading_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade,
    constraint fk_reading_device
      foreign key(device_id)
      references "device"(id)
  );

  create table "batch_reading"
  (
    id uuid primary key default uuid_generate_v4(),
    batch_id uuid not null,
    user_id uuid not null,
    reading_id uuid not null,
    created timestamptz default current_timestamp not null,
    modified timestamptz default current_timestamp not null,
    constraint fk_batch_reading_user
      foreign key(user_id)
      references "user"(id)
      on delete cascade,
    constraint fk_batch_reading_batch
      foreign key(batch_id)
      references "batch"(id)
      on delete cascade,
    constraint fk_batch_reading_reading
      foreign key(reading_id)
      references "reading"(id)
      on delete cascade
  );

commit;

do $$
declare
  admin_user uuid = uuid_generate_v4();
  admin_uuid uuid = uuid_generate_v4();
  user_uuid uuid = uuid_generate_v4();
begin
  insert into "user"
    (id, username, password)
  values
    (admin_user, 'admin', '$2a$10$5nmh/cOu.dzk05V7lfBqQua9FO6nG.aQTGTJQFB26DGMSMwp5FWxu');
  -- password = admin  

  insert into "role" (id, "name") values
    (admin_uuid, 'admin'),
    (user_uuid, 'user');

  insert into "permission" ("name") values
    ('create_user'), ('read_other_users'), ('delete_any_user'), ('delete_own_user'), ('update_own_user'), ('update_any_user'),
    ('create_batch'), ('read_batch');

  insert into "role_permission" ("role_id", "permission_id") values
  (admin_uuid, (select id from "permission" where "name" = 'create_user')),
  (admin_uuid, (select id from "permission" where "name" = 'read_other_users')),
  (admin_uuid, (select id from "permission" where "name" = 'delete_any_user')),
  (admin_uuid, (select id from "permission" where "name" = 'delete_own_user')),
  (admin_uuid, (select id from "permission" where "name" = 'update_own_user')),
  (admin_uuid, (select id from "permission" where "name" = 'update_any_user')),
  (admin_uuid, (select id from "permission" where "name" = 'create_batch')), 
  (admin_uuid, (select id from "permission" where "name" = 'read_batch')),
  
  (user_uuid, (select id from "permission" where "name" = 'update_own_user')),
  (user_uuid, (select id from "permission" where "name" = 'delete_own_user'));

  insert into "user_role" ("user_id", "role_id") values
  (admin_user, admin_uuid);
end $$;