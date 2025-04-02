package main

import (
	"Vnollx/dal/db"
	"log"
)

func main() {
	usra, _ := db.GetUserById(18)
	usrb, _ := db.GetUserById(19)
	err := db.AddFriend(usra, usrb)
	if err != nil {
		log.Println(err)
	}
}
