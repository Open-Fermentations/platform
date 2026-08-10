-- name: AddDeviceToBatch :one
insert into "batch_device" ("batch_id", "device_id") values (@batch_id, @device_id)
returning *;