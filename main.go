package main

import (
	"Simple-Distributed-Cache/cache"
	"Simple-Distributed-Cache/cluster"
	"Simple-Distributed-Cache/http"
	"Simple-Distributed-Cache/tcp"
	"flag"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)


	// 本机地址
	nodeAddr := flag.String("n", "127.0.0.1", "nodeAddr address")
	// 集群地址
	cls := flag.String("c", "", "cluster address")

	flag.Parse()

	log.Println("[nodeAddr] is", *nodeAddr)
	log.Println("[cluster] is", *cls)

	//
	node, err := cluster.New(*nodeAddr, *cls)
	if err != nil {
		panic(err)
	}

	inmemoryCache := cache.New("inmemory")
	go tcp.New(inmemoryCache, node).Listen()
	http.New(inmemoryCache, node).Listen()
}
