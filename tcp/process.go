package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
)

func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		op, err := r.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Println("close connection due to error:", err)
			}
			return
		}
		switch op {
		case 'S':
			err = s.set(conn, r)
		case 'G':
			err = s.get(conn, r)
		case 'D':
			err = s.del(conn, r)
		default:
			log.Println("close connection due to invalid operation:", op)
			return
		}
		if err != nil {
			log.Println("close connection due to error:", err)
			return
		}
	}
}
func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	// 从 conn 中读取 kv ，若 key 不属于本 node 会报错
	key, val, err := s.readKeyAndValue(r)
	if err != nil {
		return sendResponse(nil, err, conn)
	}
	// 保存 kv
	err = s.Set(key, val)
	return sendResponse(nil, err, conn)
}
func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	// 从 conn 中读取 key ，若 key 不属于本 node 会报错
	key, err := s.readKey(r)
	if err != nil {
		return sendResponse(nil, err, conn)
	}
	// 读取 kv
	val, err := s.Get(key)
	return sendResponse(val, err, conn)
}
func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	// 从 conn 中读取 key ，若 key 不属于本 node 会报错
	key, err := s.readKey(r)
	if err != nil {
		return sendResponse(nil, err, conn)
	}
	// 删除 kv
	err = s.Del(key)
	return sendResponse(nil, err, conn)
}
