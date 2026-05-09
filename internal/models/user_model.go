package models

type User struct {
	Login       string `gorm:"primaryKey"`
	Username    string
	FirstName   string
	SecondName  string
	Patronymic  string
	Stack       []Tech `gorm:"foreignKey:UserLogin;references:Login"`
	Description string
	Contacts    []Social `gorm:"foreignKey:UserLogin;references:Login"`
	ShortDesc   string
	Avatar      string
}

type CreateUserRequest struct {
	Login    string
	Email    string
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
	Stack       []Tech
	Description string
	Contacts    []Social
	ShortDesc   string
	Avatar      string
}

type PatchUserRequest struct {
	Login       string
	Username    *string
	FirstName   *string
	SecondName  *string
	Patronymic  *string
	Stack       []Tech
	Description *string
	Contacts    []Social
	ShortDesc   *string
	Avatar      *string
}

type DeleteUserRequest struct {
	Login string
}

type Social struct {
	ID        uint   `gorm:"primaryKey"`
	UserLogin string `gorm:"index"`
	Type      string
	Url       string
}

type Tech struct {
	ID        uint   `gorm:"primaryKey"`
	UserLogin string `gorm:"index"`
	Name      string
	Level     string
}
