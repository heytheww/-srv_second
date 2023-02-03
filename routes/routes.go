package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"srv_second/jwt"
	"strconv"
	"time"

	gojwt "github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

// 读写通道
var ch chan int
var ctx = context.Background()

func init() {
	ch = make(chan int, 1)
	ch <- 1
}

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

// type BuyReq struct {
// 	UserId string `json:"userid"`
// }

type BuyResp struct {
	Msg
}

func (r *Routes) HandleDfault(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte(""))
}

func (r *Routes) HandleLogin(resp http.ResponseWriter, req *http.Request) {
	user := req.FormValue("user")
	pass := req.FormValue("password")

	var v LoginResp

	j := jwt.JWT{
		HMACSecret: []byte("abc123"), // 秘钥
	}

	claims := &jwt.MyClaims{
		User: user,
	}
	claims.RegisteredClaims = gojwt.RegisteredClaims{
		ExpiresAt: gojwt.NewNumericDate(time.Now().Add(3 * time.Hour)), // 过期时间
		IssuedAt:  gojwt.NewNumericDate(time.Now()),                    // 签发时间
		NotBefore: gojwt.NewNumericDate(time.Now()),                    // 生效时间
		Issuer:    "admin",                                             // 签发人
		Subject:   "loginByPwd",                                        // 主题
		ID:        "-",
		Audience:  gojwt.ClaimStrings{"buyer"}, // 受众
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

func (r *Routes) HandleBuy(resp http.ResponseWriter, req *http.Request, rdb *redis.Client) {

	var v Msg

	// 秒杀商品ID：G18012345
	// 查看库存
	<-ch
	val, err := rdb.HGet(ctx, "GOOD_Hash_G18012345", "goodsId_count").Result()
	if err != nil {
		panic(err)
	}

	goodsId_count, err2 := strconv.Atoi(val)
	if err2 != nil {
		panic(err2)
	}

	// 消耗库存
	if goodsId_count > 0 {

		userid := req.FormValue("userid")

		if userid != "" {
			err = rdb.SAdd(ctx, "GOOD_Set_G18012345", userid).Err()
			if err != nil {
				panic(err)
			}
			err = rdb.HSet(ctx, "GOOD_Hash_G18012345", "goodsId_count", goodsId_count-1).Err()
			if err != nil {
				panic(err)
			}

			v = Msg{
				Code:    http.StatusOK,
				Massage: "秒杀成功",
			}

			vj, _ := json.Marshal(v)
			resp.Write(vj)
			return
		} else {
			v = Msg{
				Code:    http.StatusOK,
				Massage: "秒杀失败",
			}
		}

		ch <- 1

	} else {
		v = Msg{
			Code:    http.StatusOK,
			Massage: "秒杀失败",
		}

		ch <- 1
	}

	vj, _ := json.Marshal(v)
	resp.Write(vj)
}
