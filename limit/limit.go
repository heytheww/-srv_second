package limit

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

type Limit struct {
	Ctx context.Context
	lm  *rate.Limiter
}

func (l *Limit) InitLimit(r int, b int) {
	//  r 表示每r微秒产生1个令牌，b 表示令牌桶大小
	limit := rate.Every(time.Duration(r) * time.Millisecond)
	l.lm = rate.NewLimiter(limit, b)
}

func (l *Limit) ConsumeOne() {
	err := l.lm.Wait(l.Ctx)
	if err != nil {
		fmt.Println("limiter wait error: ", err)
	}
}
