package main

import (
	"flag"
)

func main() {
	vaultToken := flag.String("vaultToken", "", "Vault token")
	path := flag.String("path", "", "Absolute path to repo to lint")

	flag.Parse()

}
