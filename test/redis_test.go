package test

import (
	"fmt"
	"readygo/cache"
	"readygo/permission/permstruct"
	"strconv"
	"sync"
	"testing"
)

func init()  {
	cache.NewRedisClient(&cache.RedisConfig{
		Addr:         "127.0.0.1:6379",
	})
}


func incrWorker(id int, wg *sync.WaitGroup) {

	defer wg.Done()


		incr, _ := cache.RedisINCR("permstruct.UserStruct")
		fmt.Println(id , incr)


		var user permstruct.UserStruct


		user.Id =  strconv.Itoa(int(incr.(int64)))

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