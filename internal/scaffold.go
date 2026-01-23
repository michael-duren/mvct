package internal

import "fmt"

type ProjectConfiguration struct {
	name string
	path string
}

func ScaffoldProject(config ProjectConfiguration) {
	fmt.Println("got config: ", config)
}
