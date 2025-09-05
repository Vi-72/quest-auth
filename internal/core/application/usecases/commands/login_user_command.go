package commands

// LoginUserCommand — команда для входа пользователя
type LoginUserCommand struct {
	Email    string
	Password string
}

// LoginUserResult — результат входа
type LoginUserResult struct {
	User         UserInfo
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
}
