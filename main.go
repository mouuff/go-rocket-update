package main

import (
	"flag"
)

func main() {
	//keygenCommand := flag.NewFlagSet("keygen", flag.ExitOnError)
	signCommand := flag.NewFlagSet("sign", flag.ExitOnError)

	privateKeyPtr := signCommand.String("privateKey", "", "Private key used to sign files (Required)")
	folderPtr := signCommand.String("folder", "", "Folder to sign (Required)")

	/*
			if len(os.Args) < 2 {
		        fmt.Println("keyGen or sign command is required")
		        os.Exit(1)
			}
	*/

}
