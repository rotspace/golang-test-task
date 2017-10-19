package main

import tagc "github.com/rotspace/tagc/pkg"

func main() {
	tc := tagc.New(tagc.Port(8080))
	tc.Run()
}
