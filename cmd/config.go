package cmd

type Config struct {
	HttpPort                string
	GrpcPort                string
	DbHost                  string
	DbPort                  string
	DbUser                  string
	DbPassword              string
	DbName                  string
	DbSslMode               string
	EventGoroutineLimit     int
	JWTSecretKey            string
	JWTAccessTokenDuration  int // в минутах
	JWTRefreshTokenDuration int // в часах
}
