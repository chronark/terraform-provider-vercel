package main

import (
	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"log"
	"os"
)

func main() {
	client := vercel.New(os.Getenv("VERCEL_TOKEN"))

	err := client.Project.Delete("prj_1iprv60UGpVPyukGUp3NAJN8K4VR")
	if err != nil {
		log.Println("THE ERROR WAS HERE1")
		panic(err)
	}
}
