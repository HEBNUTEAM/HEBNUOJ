package vo

type LoginVo struct {
	NickName     string `json:"nick"`
	Email        string `json:"email"`
	Verification string `json:"verification"`
	Captcha      string `json:"captcha"`
	Password1    string `json:"pwd"`
	Password2    string `json:"pwd1"`
}
