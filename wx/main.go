package main

import (
	"readygo/wx/wxroute"
)


func main() {

	r := wxroute.NewRouter()
	r.Run(":3002")

}