package models

type ClientAuthValues struct {
	Login          string
	Password       string
	CryptoPassword string
}

type ClientRegisterResponse struct {
	Key  string
	Salt []byte
	Tag  []byte
}
