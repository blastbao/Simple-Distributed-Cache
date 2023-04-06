package cache

type Stat struct {
	Count     int64	// k-v 总数
	KeySize   int64	// k 总大小
	ValueSize int64	// v 总大小
}

func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
