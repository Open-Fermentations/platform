begin;

  create table "user"
  (
    id uuid primary key default gen_random_uuid(),
    username varchar not null,
    password varchar not null,
    active boolean default true not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null
  );

  create unique index unique_username_constraint on "user"(username);

  insert into "user"
    (username, password)
  values
    ('admin', '$2a$10$5nmh/cOu.dzk05V7lfBqQua9FO6nG.aQTGTJQFB26DGMSMwp5FWxu');
  -- password = admin

  create table "role"
  (
    id uuid primary key default gen_random_uuid(),
    "name" varchar(15) not null
  );

  create table "permission"
  (
    id uuid primary key default gen_random_uuid(),
    "name" varchar(15) not null
  );

  create table "role_permission"
  (
    id uuid primary key default gen_random_uuid(),
    role_id uuid not null,
    permission_id uuid not null,
    constraint fk_role_permission_role
      foreign key(role_id)
      references "role"(id),
    constraint fk_role_permission_permission
      foreign key(permission_id)
      references "permission"(id)
  );

  create table "brew"
  (
    id uuid primary key default gen_random_uuid(),
    "name" varchar not null,
    user_id uuid not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_brew_user
      foreign key(user_id)
      references "user"(id)
  );

  create type reading_type_enum as enum
  ('temperature', 'gravity', 'volume', 'pressure', 'ph');

  create table "device"
  (
    id uuid primary key default gen_random_uuid(),
    "name" varchar not null,
    mac_address bytea not null,
    user_id uuid not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_device_user
      foreign key(user_id)
      references "user"(id)
  );

  create unique index unique_mac_address_constraint on "device"(mac_address);

  create table "device_capability"
  (
    id uuid primary key default gen_random_uuid(),
    device_id uuid not null,
    capability reading_type_enum not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_device_capability_device
      foreign key(device_id)
      references "device"(id)
  );

  create table "reading"
  (
    id uuid primary key default gen_random_uuid(),
    reading_type reading_type_enum not null,
    "value" real not null,
    device_id uuid not null,
    user_id uuid not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_reading_user
      foreign key(user_id)
      references "user"(id),
    constraint fk_reading_device
      foreign key(device_id)
      references "device"(id)
  );

  create table "brew_reading"
  (
    id uuid primary key default gen_random_uuid(),
    brew_id uuid not null,
    user_id uuid not null,
    created timestamp default current_timestamp not null,
    modified timestamp default current_timestamp not null,
    constraint fk_brew_reading_user
      foreign key(user_id)
      references "user"(id)
  );

commit;