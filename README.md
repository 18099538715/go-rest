golang restful风格小框架，初版对参数类型、返回类型、请求方法都做了限制，目前只支持get和post方法，post方法的content-type必须为application/json
，路径参数类型暂时只支持string和int类型。抽空写的，后续会继续完善
```
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/18099538715/go-rest/webrestful"
)

type User struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

type UserController struct {
}

func add(user User) User {
	return user
}

func main() {
	handler := webrestful.Handler{}
	webrestful.Route("/aaa/{userId}", http.MethodPost, "application/json", func(userId *string, user User) User {
		b, _ := json.Marshal(user)
		fmt.Println(string(b), *userId)
		return user
	})
	webrestful.Route("/aaa/{userId}/{uesrname}", http.MethodGet, "application/json", func(userId string, userName string) string {
		return userId
	})
	http.ListenAndServe(":8000", handler)
}
```
