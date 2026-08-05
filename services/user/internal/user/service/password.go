package service

type PasswordHasher interface {
	Hash(string) (string, error)
	Compare(string, string) bool
}
