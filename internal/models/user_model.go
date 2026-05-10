package models

type User struct {
	Login       string   `json:"login" gorm:"primaryKey"`
	Username    string   `json:"username"`
	FirstName   string   `json:"fistName"`
	SecondName  string   `json:"secondName"`
	Patronymic  string   `json:"Patronimyc"`
	Stack       []Tech   `json:"stack" gorm:"foreignKey:UserLogin;references:Login"`
	Description string   `json:"description"`
	Contacts    []Social `json:"contacts" gorm:"foreignKey:UserLogin;references:Login"`
	ShortDesc   string   `json:"shortDesc"`
	Avatar      string   `json:"avatar"`
}

type CreateUserRequest struct {
	Login    string
	Username string
}

type ReadUserRequest struct {
	Login string
}

type UpdateUserRequest struct {
	Login       string
	Username    string
	FirstName   string
	SecondName  string
	Patronymic  string
	Stack       []Tech `gorm:"-"`
	Description string
	Contacts    []Social `gorm:"-"`
	ShortDesc   string
	Avatar      string
}

type PatchUserRequest struct {
	Login       string
	Username    *string
	FirstName   *string
	SecondName  *string
	Patronymic  *string
	Stack       []Tech `gorm:"-"`
	Description *string
	Contacts    []Social `gorm:"-"`
	ShortDesc   *string
	Avatar      *string
}

type DeleteUserRequest struct {
	Login string
}

type Social struct {
	ID        uint   `gorm:"primaryKey"`
	UserLogin string `gorm:"column:user_login;index"`
	Type      string `json:"type"`
	Url       string `json:"url"`
}

type Tech struct {
	ID        uint   `gorm:"primaryKey"`
	UserLogin string `gorm:"column:user_login;index"`
	Name      string `json:"name"`
	Level     string `json:"level"`
}
