package main

import (
	"readygo/wx/wxroute"
)


func main() {

	r := wxroute.NewRouter()

	//m := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist("localhost", "example2.com"),
	//	Cache:      autocert.DirCache("/var/www/.cache"),
	//}
	//
	//log.Fatal(autotls.RunWithManager(r, &m))

	r.Run(":3002")

}