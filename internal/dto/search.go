package dto

import (
	"net/http"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/route"
)

const (
	LimitDefault  int  = 50
	OffsetDefault int  = 0
	AscDefault    bool = true
)

type SearchDTO struct {
	Search  string
	Limit   int
	Offset  int
	OrderBy string
	Asc     bool
}

// Slog implements [logging.Slog].
func (s SearchDTO) Slog() []any {
	panic("unimplemented")
}

func (s *SearchDTO) FromRequest(r *http.Request) *SearchDTO {
	s.Search = route.GetStringQueryParam(r, "search", "")
	s.Limit = route.GetIntQueryParam(r, "limit", LimitDefault)
	s.Offset = route.GetIntQueryParam(r, "offset", OffsetDefault)
	s.OrderBy = route.GetStringQueryParam(r, "orderBy", "id")
	s.Asc = route.GetBoolQueryParam(r, "asc", AscDefault)

	return s
}

var _ logging.Slog = SearchDTO{}
