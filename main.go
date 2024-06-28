package main

import (
	"log"

	"github.com/toastsandwich/cli-app-with-knight/knight"
)

func Home(req *knight.Request, res *knight.Response) {
	// some logic
	username, err := req.GetParam("id")
	if err != nil {
		//
	}
	res.Write([]byte(username))

}

func main() {
	k := knight.Suitup("localhost:7000", "tcp")

	k.HandlePoint("/home", Home)
	log.Fatal(k.Serve())
}
