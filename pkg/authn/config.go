package authn

// Config is a copy of kerating/authn-go/authn.Config
// to avoid direct dependecies to authn-go
type Config struct {
	Audience           string `json:"audience" yaml:"audience"`
	Issuer             string `json:"issuer" yaml:"issuer"`
	PrivateBaseAddress string `json:"server" yaml:"server"`
	Username           string `json:"username" yaml:"username"`
	Password           string `json:"password" yaml:"password"`
}
