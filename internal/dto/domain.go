package dto

// author zhasulan
// created on 27.10.20 10:32

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Login    string `json:"login"`
	Refresh  string `json:"refresh"`
}

type NotifierRequest struct {
	Code     string `json:""`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Checksum string `json:"checksum"`
}

type CDNProxyRequest struct {
	Code     string `json:""`
	Filename string `json:"filename"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
