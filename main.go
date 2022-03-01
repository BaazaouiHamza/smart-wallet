package main

import (
	"fmt"
	"git.digitus.me/prosperus/protocol/identity"
)

func main() {
	for i := 0; i < 20; i++ {
		fmt.Println(identity.NewSecretKey().PublicKey())
	}
}
