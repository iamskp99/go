package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {

	cmd := exec.Command("grpcurl", "--plaintext", "-d", `{"name" : "Luke Skywalker","age" : "23"}`, "localhost:50051", "usermgmt.UserManagement.CreateNewUser")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("grpccurl ran successfully !!!!")
}
