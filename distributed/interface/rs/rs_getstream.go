package rs

import (
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/interface/objectstream"
	"io"
)

type RSGetStream struct {
	*decoder
}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	if len(locateInfo) + len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]			// 定位不存在的，那么保证dataServers中有补充
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		// 获取可用reader
		reader, e := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if e == nil {
			readers[i] = reader
		}
	}

	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	for i := range readers {
		if readers[i] == nil {
			// 有缺失的块， 生成对应的写入流， 因为是temp流，所以最后还会有一个Close()调用Commit(true)转正过程
			writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				return nil, e
			}
		}
	}
	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{dec}, nil

}

