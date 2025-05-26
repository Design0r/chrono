package schemas

import "chrono/db/repo"

type CreateUser struct {
	Name     string `form:"qwenameasd"     json:"name"`
	Email    string `form:"qweemailasd"    json:"email"`
	Password string `form:"qwepasswordasd" json:"password"`
}

type PatchUser struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Color string `form:"color"`
}

type Login struct {
	Email    string `form:"qweemailasd"`
	Password string `form:"qwepasswordasd"`
}

type VacUser struct {
	PlannedVacation int
	repo.User
}

type DebugUsers struct {
	Users []repo.CreateUserParams `json:"users"`
}

type Honeypot struct {
	Name     string `form:"name"     json:"name"`
	Email    string `form:"email"    json:"email"`
	Password string `form:"password" json:"password"`
}

func (h Honeypot) IsFilled() bool {
	return h.Name != "" || h.Email != "" || h.Password != ""
}
