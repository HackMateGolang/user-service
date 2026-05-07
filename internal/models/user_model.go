package models

type User struct {
	Login       string `gorm:"primaryKey"`
	Username    string
	FirstName   string
	SecondName  string
	Patronymic  string
	Stack       map[string]struct{}
	Description string
	Сontacts    map[string]string
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
	Stack       map[string]struct{}
	Description string
	Сontacts    map[string]string
	ShortDesc   string
	Avatar      string
}

type PatchUserRequest struct {
	Login       string
	Username    *string
	FirstName   *string
	SecondName  *string
	Patronymic  *string
	Stack       *map[string]struct{}
	Description *string
	Сontacts    *map[string]string
	ShortDesc   *string
	Avatar      *string
}

type DeleteUserRequest struct {
	Login string
}
