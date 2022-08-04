package validate

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	//正则判断是否合法
	if ok,_ := regexp.MatchString(`^(1[3-8])\d{9}$`, mobile); !ok {
		return false
	}
	return true
}