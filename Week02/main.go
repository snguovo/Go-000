package main

import (
	"log"

	"github.com/snguovo/Go-000/Week02/tree/main/service"
)

func main() {
	srv := service.NewService()
	_, err := srv.FindUserByID(1)
	log.Println(err)
	_, err = srv.MustFindUserByID(1)
	log.Println(err)
}
