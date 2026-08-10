-- name: AddDeviceToBatch :one
insert into "batch_device" ("batch_id", "device_id") values (@batch_id, @device_id)
returning *;

-- name: RemoveDeviceFromBatch :exec
delete from "batch_device"
where "batch_id" = @batch_id and "device_id" = @device_id;