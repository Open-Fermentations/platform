begin;

alter table "batch_reading" drop constraint fk_batch_reading_user;
drop table if exists "batch_reading";

alter table "reading" drop constraint fk_reading_device;
alter table "reading" drop constraint fk_reading_user;
drop table if exists "reading";

alter table "device_capability" drop constraint fk_device_capability_device;
drop table if exists "device_capability";

drop index unique_mac_address_constraint;
alter table "device" drop constraint fk_device_user;
drop table if exists "device";

drop type if exists reading_type_enum;

alter table "batch" drop constraint fk_batch_user;
drop table if exists "batch";

alter table "role_permission" drop constraint fk_role_permission_permission;
alter table "role_permission" drop constraint fk_role_permission_role;
drop table if exists "role_permission";

drop table if exists "permission";

drop table if exists "role";

drop index unique_username_constraint;
drop table if exists "user";

commit;