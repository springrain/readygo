package main

import (
	"os"
	"readygo/wx/wxroute"
)


 var systemPath, _ = os.Getwd()

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