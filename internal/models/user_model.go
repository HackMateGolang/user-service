package models

type User struct {
	Login       string
	Username    string
	FirstName   string
	SecondName  string
	Patronymic  string
	Stack       map[string]struct{}
	Description string
	contacts    map[string]string
	ShortDesc   string
	Avatar      string
}

type CreateUserRequest struct {
	login    string
	email    string
	username string
}

type ReadUserRequest struct {
	login string
}

type UpdateUserRequest struct {
	Login       string
	Username    string
	FirstName   string
	SecondName  string
	Patronymic  string
	Stack       map[string]struct{}
	Description string
	contacts    map[string]string
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
	contacts    *map[string]string
	ShortDesc   *string
	Avatar      *string
}

type DeleteUserRequest struct {
	login string
}
