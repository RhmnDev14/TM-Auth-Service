package main

import (
	"auth-service/internal"
)

func main() {
	internal.NewServer().Run()
}
