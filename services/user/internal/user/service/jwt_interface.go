package service

type JwtService interface {
	Generate(payload any) (string, error)
	GenerateRefresh(payload any) (string, error)
}
