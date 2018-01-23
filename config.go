package servicecomb

var DEFAULT_SC_ADDR = "http://127.0.0.1:30100"

type ServiceComb_Config struct {
	SCAddr string `json:"sc_addr"`
}

var DEFAULT_SC_CONFIG = ServiceComb_Config{
	SCAddr: DEFAULT_SC_ADDR,
}
