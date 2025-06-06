package params

type BaseQueryParams struct {
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
}

type ServiceQueryParams struct {
	BaseQueryParams
	UserID     string `query:"user_id"`
	Status     string `query:"status"`
	BookedDate string `query:"booked_date"`
	BranchID   string `query:"branch_id"`
	CategoryID string `query:"category_id"`
	BookedTime string `query:"booked_time"`
}

func NewServiceQueryParams() *ServiceQueryParams {
	return &ServiceQueryParams{
		BaseQueryParams: BaseQueryParams{
			Limit:  10,
			Offset: 0,
		},
	}
}
