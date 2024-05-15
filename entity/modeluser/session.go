package modeluser

// UserSession 存储的用户session信息
type UserSession struct {
	ID          string `json:"id"`
	UserAccount string `json:"user_account"`
	UserRole    string `json:"user_role"`
}

func UserToUserSession(u User) UserSession {
	var session UserSession
	session.ID = u.ID
	session.UserRole = u.UserRole
	session.UserAccount = u.UserAccount

	return session
}
