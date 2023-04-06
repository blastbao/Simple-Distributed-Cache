package http

import (
	"Simple-Distributed-Cache/cache"
	"Simple-Distributed-Cache/cluster"
	"net/http"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	// 读写 cache
	http.Handle("/cache/", s.cacheHandler())
	// 获取 cache 的 k-v 总数、key 大小、value 大小
	http.Handle("/status", s.statusHandler())
	// 获取 cluster 成员列表
	http.Handle("/cluster", s.clusterHandler())
	// 遍历本 node 的 cache 中的 kvs ，将非本 node 的 kvs 转发到其它 nodes
	http.Handle("/rebalance", s.rebalanceHandler())
	http.ListenAndServe(s.Addr()+":12345", nil)
}

func New(cache cache.Cache, node cluster.Node) *Server {
	return &Server{cache, node}
}
