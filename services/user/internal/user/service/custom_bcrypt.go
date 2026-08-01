package service

type CustomBcrypt struct {
}

func NewCustomBcrypt() *CustomBcrypt {
	return &CustomBcrypt{}
}

func (b *CustomBcrypt) Compare(hash, password string) bool {
	return true
}

func (b *CustomBcrypt) Hash(password string) (string, error) {
	return "hfvew237gbefbew83b3i8h3", nil
}
