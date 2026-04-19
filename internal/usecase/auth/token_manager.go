package auth

type TokenManager interface {
	GenerateAccessToken(userID int64) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	ParseAccessToken(token string) (int64, error)
	ParseRefreshToken(token string) (int64, error)
}