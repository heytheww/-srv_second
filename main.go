package main

import (
	"context"
	"net/http"
	"srv_second/limit"
	"srv_second/routes"

	"srv_second/redis"

	goredis "github.com/redis/go-redis/v9"
)

// 自定义路由
type MyMux struct {
	RDB     *goredis.Client
	Limiter *limit.Limit
}

func (p *MyMux) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	r := routes.Routes{}

	if req.URL.Path == "/" {
		resp.Write([]byte(""))
		r.HandleDfault(resp, req)
		return
	}
	if req.URL.Path == "/login" {
		r.HandleLogin(resp, req)
		return
	}
	if req.URL.Path == "/buy" {
		// 限流:1个令牌/秒，桶容量100
		p.Limiter.ConsumeOne()

		r.HandleBuy(resp, req, p.RDB)
		return
	}
	http.NotFound(resp, req)
}

func main() {

	var ctx = context.Background()
	rd := redis.Redis{
		Ctx:      ctx,
		Addr:     "43.136.78.171:30120",
		Password: "test001",
		DB:       0,
	}
	var err error
	rd.Init()
	rdb := rd.GetDB()

	// 初始化商品 G18012345 信息，这步应在其他模块，比如商品运维模块完成
	err = rdb.HSet(ctx, "GOOD_Hash_G18012345", "goodsId_count", 0).Err()
	if err != nil {
		panic(err)
	}

	limiter := limit.Limit{Ctx: context.Background()}
	limit := 1000 // 1s
	limiter.InitLimit(limit, 10)

	mux := &MyMux{RDB: rdb, Limiter: &limiter}
	http.ListenAndServe(":1234", mux)
}
