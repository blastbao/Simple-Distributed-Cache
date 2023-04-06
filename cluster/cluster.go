package cluster

import (
	"io/ioutil"
	"time"

	"github.com/hashicorp/memberlist"
	"stathat.com/c/consistent"
)

type Node interface {
	ShouldProcess(key string) (string, bool) // 获取 key 归属的 node 地址
	Members() []string //该函数被consistent实现
	Addr() string
}
type node struct {
	*consistent.Consistent // 集群 nodes 映射的 hash 环，定时(1s)更新
	addr string	// 本机地址
}

func New(addr, cluster string) (Node, error) {
	//创建gossip新节点的config
	config := memberlist.DefaultLANConfig()
	config.Name = addr
	config.BindAddr = addr
	config.LogOutput = ioutil.Discard
	//创建新节点
	mbl, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}
	if cluster == "" {
		cluster = addr
	}
	existing := []string{cluster}
	//连接到集群
	_, err = mbl.Join(existing)
	if err != nil {
		return nil, err
	}
	//创建一致性哈希的节点实例
	circle := consistent.New()
	//设置虚拟节点数量
	circle.NumberOfReplicas = 256
	go func() {
		for {
			m := mbl.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			//每隔1s将集群节点列表m更新到circle中
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()
	return &node{circle, addr}, nil

}

// ShouldProcess 获取 key 归属的 node 地址，以及其是否属于本 node
func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}

func (n *node) Addr() string {
	return n.addr
}
