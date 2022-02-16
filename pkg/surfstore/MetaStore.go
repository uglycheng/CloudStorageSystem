package surfstore

import (
	context "context"
	"sync"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type MetaStore struct {
	FileMetaMap    map[string]*FileMetaData
	BlockStoreAddr string
	metaLock       sync.Mutex
	UnimplementedMetaStoreServer
}

func (m *MetaStore) GetFileInfoMap(ctx context.Context, _ *emptypb.Empty) (*FileInfoMap, error) {
	copyMap := make(map[string]*FileMetaData)
	m.metaLock.Lock()
	for k, v := range m.FileMetaMap {
		copyMap[k] = v
	}
	m.metaLock.Unlock()

	return &FileInfoMap{FileInfoMap: copyMap}, nil
}

func (m *MetaStore) UpdateFile(ctx context.Context, fileMetaData *FileMetaData) (*Version, error) {
	m.metaLock.Lock()
	if _, ok := m.FileMetaMap[fileMetaData.Filename]; ok {
		if fileMetaData.Version == (m.FileMetaMap[fileMetaData.Filename].Version + 1) {
			m.FileMetaMap[fileMetaData.Filename] = fileMetaData
			m.metaLock.Unlock()
			return &Version{Version: fileMetaData.Version}, nil
		} else {
			m.metaLock.Unlock()
			return &Version{Version: -1}, nil
		}
	} else {
		m.FileMetaMap[fileMetaData.Filename] = fileMetaData
		m.metaLock.Unlock()
		return &Version{Version: 1}, nil
	}
}

func (m *MetaStore) GetBlockStoreAddr(ctx context.Context, _ *emptypb.Empty) (*BlockStoreAddr, error) {
	return &BlockStoreAddr{Addr: m.BlockStoreAddr}, nil
}

// This line guarantees all method for MetaStore are implemented
var _ MetaStoreInterface = new(MetaStore)

func NewMetaStore(blockStoreAddr string) *MetaStore {
	return &MetaStore{
		FileMetaMap:    map[string]*FileMetaData{},
		BlockStoreAddr: blockStoreAddr,
	}
}
