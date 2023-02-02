package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"srv_second/jwt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v4"
)

type Routes struct {
}

type Msg struct {
	Code    int    `json:"code"`
	Massage string `json:"massage"`
}

type LoginResp struct {
	Msg
	Token string `json:"token"`
}

type BuyResp struct {
	Msg
}

func (r *Routes) HandleLogin(resp http.ResponseWriter, req *http.Request) {
	user := req.FormValue("user")
	pass := req.FormValue("password")

	var v LoginResp

	j := jwt.JWT{
		HMACSecret: []byte("abc123"), // 秘钥
	}

	// 偷懒格式化
	claims := &jwt.MyClaims{
		user,
		gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(3 * time.Hour)), // 过期时间
			IssuedAt:  gojwt.NewNumericDate(time.Now()),                    // 签发时间
			NotBefore: gojwt.NewNumericDate(time.Now()),                    // 生效时间
			Issuer:    "admin",                                             // 签发人
			Subject:   "loginByPwd",                                        // 主题
			ID:        "-",
			Audience:  gojwt.ClaimStrings{"buyer"}, // 受众
		},
	}

	token, err := j.Sign(claims)
	if err != nil {
		panic(err)
	}

	if user == "admin" && pass == "123" {
		msg := Msg{
			Code:    http.StatusOK,
			Massage: "登录成功",
		}

		v = LoginResp{
			Token: token,
		}
		v.Msg = msg
	} else {
		msg := Msg{
			Code:    http.StatusOK,
			Massage: "登录失败",
		}

		v = LoginResp{
			Token: "",
		}
		v.Msg = msg
	}

	vj, _ := json.Marshal(v)
	resp.Write(vj)
}

func (r *Routes) HandleBuy(resp http.ResponseWriter, req *http.Request) {

	fmt.Println("执行业务逻辑")
	v := Msg{
		Code:    http.StatusOK,
		Massage: "",
	}
	vj, _ := json.Marshal(v)

	resp.Write(vj)
}
