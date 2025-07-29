package msgs

type ErrorMsg error

type SelectAuthMsg struct {
	Address  string
	Selected string
}

type FormSubmitMsg struct {
	Values     map[string]string
	ButtonName string
}

type CredentialsBytesMsg struct {
	Key  string
	Salt []byte
	Tag  []byte
}

type UpdateKeySuccessMsg string

type CryptoPassValid struct{}
