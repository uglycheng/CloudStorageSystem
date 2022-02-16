package surfstore

import (
	context "context"
	"fmt"
	"sync"
)

type BlockStore struct {
	BlockMap  map[string]*Block
	blockLock sync.Mutex
	UnimplementedBlockStoreServer
}

func (bs *BlockStore) GetBlock(ctx context.Context, blockHash *BlockHash) (*Block, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	bs.blockLock.Lock()
	block := bs.BlockMap[blockHash.Hash]
	bs.blockLock.Unlock()
	return block, nil
}

func (bs *BlockStore) PutBlock(ctx context.Context, block *Block) (*Success, error) {
	if ctx.Err() != nil {
		fmt.Println(ctx.Err())
		return nil, ctx.Err()
	}
	hash := GetBlockHashString(block.BlockData)
	bs.blockLock.Lock()
	bs.BlockMap[hash] = block
	bs.blockLock.Unlock()
	return &Success{Flag: true}, nil
}

// Given a list of hashes “in”, returns a list containing the
// subset of in that are stored in the key-value store
func (bs *BlockStore) HasBlocks(ctx context.Context, blockHashesIn *BlockHashes) (*BlockHashes, error) {
	hashIn := make([]string, 0)
	bs.blockLock.Lock()
	for _, hash := range blockHashesIn.Hashes {
		if _, ok := bs.BlockMap[hash]; ok {
			hashIn = append(hashIn, hash)
		}
	}
	bs.blockLock.Unlock()
	return &BlockHashes{Hashes: hashIn}, nil
}

// This line guarantees all method for BlockStore are implemented
var _ BlockStoreInterface = new(BlockStore)

func NewBlockStore() *BlockStore {
	return &BlockStore{
		BlockMap: map[string]*Block{},
	}
}
