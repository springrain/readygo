package test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"readygo/cache"
	"readygo/permission/permstruct"

	"gitee.com/chunanyong/zorm/decimal"
)

type Uu struct {
	Name    string
	Teacher string
	Amount  int
	o       decimal.Decimal

	E decimal.Decimal
	d string
}

func init() {
	ctx := context.Background()
	cache.NewRedisClient(ctx, &cache.RedisConfig{
		Addr: "127.0.0.1:6379",
	})

	cache.NewRedisCacheManager()
}

func incrWorker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()
	incr, _ := cache.RedisINCR(ctx, "permstruct.UserStruct")

	if incr == nil {
		incr, _ = cache.RedisINCR(ctx, "permstruct.UserStruct")
	}
	fmt.Println(id, incr)

	var user permstruct.UserStruct

	user.Id = strconv.Itoa(int(incr.(int64)))

	fmt.Println(user)
}

func TestINCR(t *testing.T) {
	var wg sync.WaitGroup

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go incrWorker(i, &wg)
	}

	wg.Wait()
}

func TestDemo(t *testing.T) {
	ctx := context.Background()

	incr, err := cache.RedisINCR(ctx, "permstruct.UserStruct")

	fmt.Println(incr, err)
}

func TestStruct(t *testing.T) {
	var u Uu
	u.Amount = 2
	u.Name = "fef"
	u.Teacher = "yyy"
	u.o = decimal.NewFromFloat(2.33)

	u.E = decimal.NewFromFloat(2.333)

	u.d = "ff"
	ctx := context.Background()
	err := cache.PutToCache(ctx, "Uu.table", "dd", &u)

	fmt.Println(err, u)

	r := &Uu{}

	cache.GetFromCache(ctx, "Uu.table", "dd", r)

	fmt.Println(r)
}
