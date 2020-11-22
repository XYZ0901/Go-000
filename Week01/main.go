package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("hello,world!")
}

// service
type service struct{}

// Subset 子集算法 backends: 服务节点  clientId: 客户端id  subsetSize: 子集大小
func Subset(backends []service, clientId, subsetSize int) []service {
	subsetCount := len(backends) / subsetSize

	round := clientId / subsetCount
	rand.Seed(int64(round))
	rand.Shuffle(len(backends), func(i, j int) {
		backends[i], backends[j] = backends[j], backends[i]
	})

	subsetId := clientId % subsetCount

	start := subsetId * subsetSize
	return backends[start : start+subsetSize]
}
