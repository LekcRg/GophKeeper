package models

type ClientRegisterValues struct {
	Login          string
	Password       string
	CryptoPassword string
}

type ClientRegisterResponse struct {
	Key  string
	Salt []byte
	Tag  []byte
}
