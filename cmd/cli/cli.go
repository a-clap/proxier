package main

import (
	"encoding/json"
	"fmt"
)

type TestJson struct {
	User     string  `json:"user"`
	Password *string `json:"password"`
}

func main() {
	tmp := `{
			"user": "adam"
			}`
	tst := &TestJson{}
	if err := json.Unmarshal([]byte(tmp), tst); err != nil {
		panic(err)
	}
	fmt.Println(tst)

}
