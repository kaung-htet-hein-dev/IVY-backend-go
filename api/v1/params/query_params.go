package params

type BaseQueryParams struct {
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
}

// service
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

// user

type UserQueryParams struct {
	BaseQueryParams
	Name        string `query:"name"`
	Email       string `query:"email"`
	Password    string `query:"password"`
	PhoneNumber string `query:"phone_number"`
	Role        string `query:"role"`
}

func NewUserQueryParams() *UserQueryParams {
	return &UserQueryParams{
		BaseQueryParams: BaseQueryParams{
			Limit:  10,
			Offset: 0,
		},
	}
}

// branch

type BranchQueryParams struct {
	BaseQueryParams
	Name        string `query:"name"`
	Location    string `query:"location"`
	PhoneNumber string `query:"phone_number"`
	ServiceID   string `query:"service_id"`
}

func NewBranchQueryParams() *BranchQueryParams {
	return &BranchQueryParams{
		BaseQueryParams: BaseQueryParams{
			Limit:  10,
			Offset: 0,
		},
	}
}
