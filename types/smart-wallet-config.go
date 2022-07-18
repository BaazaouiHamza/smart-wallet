package types

type SmartWalletConfig struct {
	DB struct {
		Host       string `json:"host"`
		Name       string `json:"name"`
		Pass       string `json:"pass"`
		Port       int    `json:"port"`
		User       string `json:"user"`
		SSLEnabled bool   `json:"sslEnabled"`
	} `json:"db"`
	Address                   string `json:"address"`
	Port                      int    `json:"port"`
	ConsulAddress             string `json:"consulAddress"`
	NsqLookupAddress          string `json:"nsqLookUpAddress"`
	NsqLookupAddressProsperus string `json:"nsqLookUpAddressProsperus"`
	ClientCredentials         struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	} `json:"clientCredentials"`
}
