package todo

import "time"

type SortingOrder string
type SortingType string

var SortingOrders = struct {
	Asc  SortingOrder
	Desc SortingOrder
}{
	Asc:  "ASC",
	Desc: "DESC",
}

func (so *SortingOrder) String() string {
	return string(*so)
}

var SortingTypes = struct {
	CreatedAt   SortingType
	UpdatedAt   SortingType
	PublishedAt SortingType
	ID          SortingType
	Body        SortingType
	Status      SortingType
}{
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	PublishedAt: "published_at",
	ID:          "id",
	Body:        "body",
	Status:      "status",
}

func (st *SortingType) String() string {
	return string(*st)
}

var ModelTypes = struct {
	Report      int
	Message     int
	ReportReply int
}{
	Message:     1,
	Report:      2,
	ReportReply: 3,
}

type CursorPagingParam struct {
	Token        *string      `json:"token"`
	Size         int          `json:"size" validate:"gt=0"`
	HasCursor    bool         `json:"has_cursor"`
	SortingOrder SortingOrder `json:"-"`
	LastDataAt   *time.Time   `json:"-"`
}

type CursorPagingTime struct {
	CursorPaging time.Time
}
