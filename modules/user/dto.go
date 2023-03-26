package user

type InitUserDTO struct {
	Name     string `json:"name" form:"name" validate:"omitempty"`
	Username string `json:"username" form:"username" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	//check with authentication process
	Password string `json:"password" form:"password" validate:"required"`
	PhotoUrl string `json:"photoUrl" form:"photoUrl" validate:"omitempty,url"`
	Gender   uint8  `json:"gender" form:"gender" validate:"omitempty,number,gte=1,lte=3"`
	Birthday string `json:"birthday" form:"birthday" validate:"omitempty,datetime=2006-01-02"`
}

type UpdateUserDTO struct {
	Username string `json:"username" form:"username" validate:"omitempty"`
	Name     string `json:"name" form:"name" validate:"omitempty"`
	Uid      int    `json:"uid" form:"uid" validate:"required"`
	Email    string `json:"email" form:"email" validate:"omitempty"`
	Gender   uint8  `json:"gender" form:"gender" validate:"omitempty,number,gte=1,lte=3"`
	Password string `json:"password" form:"password" validate:"omitempty"`
	PhotoUrl string `json:"photoUrl" form:"photoUrl" validate:"omitempty,url"`
	// Note: use RFC3339 Format for the date string
	Birthday string `json:"birthday" form:"birthday" validate:"omitempty,datetime=2006-01-02"`
}
