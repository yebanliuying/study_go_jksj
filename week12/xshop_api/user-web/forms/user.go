package forms

type PasswordLoginForm struct {
	Mobile    string `json:"mobile" binding:"required,mobile"`
	Password  string `json:"password" binding:"required,min=6,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Mobile   string `json:"mobile" binding:"required,mobile"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Code     string `json:"code" binding:"required,min=4,max=4"`
}
