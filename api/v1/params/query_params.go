package params

type BaseQueryParams struct {
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
}

// booking
type BookingQueryParams struct {
	BaseQueryParams
	UserID     string `query:"user_id"`
	Status     string `query:"status"`
	BookedDate string `query:"booked_date"`
	BranchID   string `query:"branch_id"`
	CategoryID string `query:"category_id"`
	BookedTime string `query:"booked_time"`
}

func NewBookingQueryParams() *BookingQueryParams {
	return &BookingQueryParams{
		BaseQueryParams: BaseQueryParams{
			Limit:  10,
			Offset: 0,
		},
	}
}

// user

type UserQueryParams struct {
	BaseQueryParams
	ID          string `query:"id"`
	Name        string `query:"name"`
	Email       string `query:"email"`
	PhoneNumber string `query:"phone_number"`
	Role        string `query:"role"`
	Verified    bool   `query:"verified"`
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
	IsActive    *bool  `query:"is_active"`
}

func NewBranchQueryParams() *BranchQueryParams {
	return &BranchQueryParams{
		BaseQueryParams: BaseQueryParams{
			Limit:  10,
			Offset: 0,
		},
	}
}

// service

type ServiceQueryParams struct {
	BaseQueryParams
	Name           string `query:"name"`
	DurationMinute int    `query:"duration_minute"`
	Price          int    `query:"price"`
	CategoryID     string `query:"category_id"`
	BranchID       string `query:"branch_id"`
}

func NewServiceQueryParams() *ServiceQueryParams {
	return &ServiceQueryParams{
		BaseQueryParams: BaseQueryParams{
			Limit:  10,
			Offset: 0,
		},
	}
}

// category

type CategoryQueryParams struct {
	BaseQueryParams
	Name string `query:"name"`
}

func NewCategoryQueryParams() *CategoryQueryParams {
	return &CategoryQueryParams{
		BaseQueryParams: BaseQueryParams{
			Limit:  10,
			Offset: 0,
		},
	}
}
