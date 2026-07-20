package database

import (
	"context"
	"fmt"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/model"
)

const searchDevices = `
select id, name, mac_address, user_id, created, modified, count(id) over() as total from "device"
where "name" like $1::text
order by "%s" %s
limit $2::integer
offset $3:: integer
`

type SearchDevicesParams struct {
	Name   string
	Order  string
	Asc    bool
	Offset int32
	Limit  int32
}

func (d service) SearchDevices(ctx context.Context, arg SearchDevicesParams) (*dto.PageDTO[model.Device], error) {
	var direction string
	if arg.Asc {
		direction = "asc"
	} else {
		direction = "desc"
	}
	rows, err := d.dbpool.Query(ctx, fmt.Sprintf(searchDevices, arg.Order, direction),
		arg.Name,
		arg.Limit,
		arg.Offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var total int64
	var items []model.Device
	for rows.Next() {
		var i model.Device
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MacAddress,
			&i.UserID,
			&i.Created,
			&i.Modified,
			&total,
		); err != nil {
			return nil, err
		}

		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &dto.PageDTO[model.Device]{
		Limit:  int(arg.Limit),
		Offset: int(arg.Offset),
		Total:  int(total),
		Data:   items,
	}, nil
}
