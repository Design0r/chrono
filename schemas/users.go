package schemas

import "chrono/db/repo"

type CreateUser struct {
	Name     string `form:"name"     json:"name"`
	Email    string `form:"email"    json:"email"`
	Password string `form:"password" json:"password"`
}

type PatchUser struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Color string `form:"color"`
}

type Login struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type VacUser struct {
	PlannedVacation int
	repo.User
}

type DebugUsers struct {
	Users []repo.CreateUserParams `json:"users"`
}
