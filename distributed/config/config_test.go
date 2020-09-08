package config

import (
	"fmt"
	"testing"
)

func Test_Init(t *testing.T) {
	fmt.Println(Pub)
	fmt.Printf("%p\n", Pub)
	fmt.Println(ServerInf)
	fmt.Println(ServerData)
}
