package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

//MarshalJSON 重写方法
func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02 15:04:05"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32    `json:id`
	Nickname string   `json:nickname`
	Mobile   string   `json:mobile`
	Gender   uint32   `json:gender`
	Birthday JsonTime `json:birthday`
}
