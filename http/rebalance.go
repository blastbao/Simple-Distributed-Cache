package http

import (
	"bytes"
	"net/http"
)

type rebalanceHandler struct {
	*Server
}

func (h *rebalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	go h.rebalance()
}

func (h *rebalanceHandler) rebalance() {
	s := h.NewScanner()
	defer s.Close()
	client := &http.Client{}

	// 遍历本 node 的 cache 中的 kvs ，将非本 node 的 kvs 转发到其它 nodes
	for s.Scan() {
		k := s.Key()
		// 获取 key 归属的 node 地址，以及其是否属于本 node
		redirectAddr, ok := h.ShouldProcess(k)
		// 若 key 当前不归属于本 node ，则转存到其它 node
		if !ok {
			// 通过 http 请求将 kv 转发到目标 node
			r, _ := http.NewRequest(http.MethodPut, "http://"+redirectAddr+":12345/cache/"+k, bytes.NewReader(s.Value()))
			client.Do(r)
			h.Del(k)
		}
	}
}

func (s *Server) rebalanceHandler() http.Handler {
	return &rebalanceHandler{s}
}
