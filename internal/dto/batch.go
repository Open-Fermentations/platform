package dto

import (
	"log/slog"
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/model"
	"time"

	"github.com/google/uuid"
)

type CreateBatchDTO struct {
	Name string `json:"name"`
}

func (c *CreateBatchDTO) ToCreateBatchParams(id uuid.UUID) *sqlc.CreateBatchParams {
	return &sqlc.CreateBatchParams{
		Name:   c.Name,
		UserID: id,
	}
}

// Slog implements [logging.Slog].
func (c CreateBatchDTO) Slog() []any {
	return []any{slog.Group("CreateBatchDto", slog.String("name", c.Name))}
}

var _ logging.Slog = CreateBatchDTO{}

type BatchDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	UserID   uuid.UUID `json:"userId"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// Slog implements [logging.Slog].
func (d BatchDTO) Slog() []any {
	return []any{
		slog.Group("batch",
			slog.String("id", d.ID.String()),
			slog.String("name", d.Name),
			slog.String("userId", d.UserID.String()),
			slog.Time("created", d.Created),
			slog.Time("modified", d.Modified),
		)}
}

func (d *BatchDTO) FromModel(m model.Batch) *BatchDTO {
	d.ID = m.ID
	d.Name = m.Name
	d.UserID = m.UserID
	d.Created = m.Created
	d.Modified = m.Modified

	return d
}

var _ logging.Slog = BatchDTO{}

type UpdateBatchDTO struct {
	Name string `json:"name"`
}

// Slog implements [logging.Slog].
func (u UpdateBatchDTO) Slog() []any {
	return []any{slog.Group("update_batch_dto", slog.String("name", u.Name))}
}

var _ logging.Slog = UpdateBatchDTO{}
