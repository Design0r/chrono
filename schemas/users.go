package schemas

type CreateUser struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Vacation int    `form:"vacation"`
}

type Login struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}
