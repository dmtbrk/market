package user

type User struct {
	ID      string
	Name    string
	Balance int64
}

type UpdateRequest struct {
	Balance *int64
}
