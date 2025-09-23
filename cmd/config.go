package cmd

type Config struct {
	HTTPPort                string
	GrpcPort                string
	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPassword              string
	DBName                  string
	DBSslMode               string
	EventGoroutineLimit     int
	JWTSecretKey            string
	JWTAccessTokenDuration  int // в минутах
	JWTRefreshTokenDuration int // в часах
}
