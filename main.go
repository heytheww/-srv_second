package main

import (
	"context"
	"net/http"
	"srv_second/limit"
	"srv_second/routes"
)

// https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch5-web/ch5-03-middleware.html
// 中间件
func limitMiddleware(next http.Handler, lm *limit.Limit) func(http.ResponseWriter, *http.Request) {
	// 限流准备:1个令牌/秒，桶容量100

	return func(w http.ResponseWriter, r *http.Request) {
		lm.ConsumeOne()
		// next handler
		next.ServeHTTP(w, r)
	}
}

func main() {

	lm := limit.Limit{Ctx: context.Background()}
	limit := 10000 // 10s
	lm.InitLimit(limit, 10)

	r := routes.Routes{}
	http.HandleFunc("/login", r.HandleLogin)
	http.HandleFunc("/buy", limitMiddleware(http.HandlerFunc((r.HandleBuy)), &lm))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})
	http.ListenAndServe(":1234", nil)
}
