package main

import (
	"Week02/service"
	"log"
	"net/http"
)

// TODO: 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
// TODO: 需要Wrap上抛 因为sql.ErrNoRows是基础库抛出的err 所以需要上抛
func main() {
	http.HandleFunc("/authentication", service.Authentication)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
