package commands

// RegisterUserCommand — команда для регистрации пользователя
type RegisterUserCommand struct {
	Email    string
	Phone    string
	Name     string
	Password string
}

// RegisterUserResult — результат регистрации
type RegisterUserResult struct {
	User         UserInfo
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
}
