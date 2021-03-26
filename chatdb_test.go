package main

import (
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {
	DB_open()
	defer DB_close()

	fmt.Println(DeletePostOnid(10))
}
