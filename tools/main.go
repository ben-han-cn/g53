package main

import (
	"fmt"
	"github.com/ben-han-cn/g53"
)

func main() {
	c, _ := g53.NewName("c", false)
	fmt.Printf("name is %v\n", c.String(true))
}
