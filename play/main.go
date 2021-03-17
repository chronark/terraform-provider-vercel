package main

import (
	"log"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
)

func main() {
	client := vercel.New("wsByP9ptGqn7snGvvY00aDzn")

	err := client.Project.Delete("prj_1iprv60UGpVPyukGUp3NAJN8K4VR")
	if err != nil {
		log.Println("THE ERROR WAS HERE1")
		panic(err)
	}
}
