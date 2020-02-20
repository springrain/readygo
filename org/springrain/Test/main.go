package main

import (
	"fmt"
	"strings"
)

func test(a []string) []string {
	return a
}
func main() {
	str := "select * from t_user where a=? and b=?"
	values := strings.Split(str, "?")
	fmt.Println(len(values), values[0], ":", values[1], ":", len(values[2]))
}
