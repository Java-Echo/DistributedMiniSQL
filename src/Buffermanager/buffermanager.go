package Buffermanager

import "sync"
import "os"
import "fmt"

const GoFlushNum = 5 // 最多使用5个协程处理flush
const blockSize = 8192

var blockBuffer *LRUCache
var connector = "*"

var fileNamePos2Int map[nameAndPos]int

var query2IntLock sync.Mutex //
var blockNumLock sync.Mutex  //互斥锁，不加读写锁

var posNum = 0

type nameAndPos struct {
	filename string
	blockid  uint16
}

// InitBuffer
func InitBuffer() {
	blockBuffer = NewLRUCache()
	fileNamePos2Int = make(map[nameAndPos]int, InitSize*4)
	posNum = 0
}

//BlockRead 读byte，不检查blockid和filename， 同时加互斥锁（为什么不加读写锁？
func BlockRead(filename string, block_id uint16) (*block, error) {
	t := Query2Int(nameAndPos{filename: filename, blockid: block_id})
	ok, block_t := blockBuffer.findBlock(t)
	if ok {
		block_t.Lock()
		return block_t, nil
	}
	newBlock := block{
		blockid:  block_id,
		filename: filename,
		Data:     make([]byte, blockSize),
	}
	err := newBlock.read()
	if err != nil {
		return nil, err
	}
	blockPtr := blockBuffer.insertBlock(&newBlock, t)
	blockPtr.Lock()
	return blockPtr, nil
}

// GetBlockNumber 返回当前块数 大小为BlockSize 的倍数, 当前文件大小为 BlockSize * BlockNumber
func GetBlockNumber(filename string) (uint16, error) {
	blockNumLock.Lock()
	defer blockNumLock.Unlock()
	return findBlockNumber(filename)
}

// NewBlock 返回block_id指新的块在文件中的block id
func NewBlock(filename string) (uint16, error) {
	blockNumLock.Lock()
	defer blockNumLock.Unlock()
	block_id, err := findBlockNumber(filename)
	if err != nil {
		return 0, err
	}
	newBlock := block{
		blockid:  block_id,
		filename: filename,
		Data:     make([]byte, blockSize),
	}
	newBlock.setDirty()
	newBlock.flush()
	t := Query2Int(nameAndPos{filename: filename, blockid: block_id})
	blockBuffer.insertBlock(&newBlock, t)
	return block_id, nil
}

func findBlockNumber(filename string) (uint16, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	fmt.Println("the size is ", fileInfo.Size())
	tmp := fileInfo.Size() / blockSize
	return uint16(tmp), nil
}

// BlockFlushALL 刷新所有块，使用协程处理
func BlockFlushAll() (bool, error) {
	blockBuffer.Lock()
	defer blockBuffer.Unlock()
	flushChannel := make(chan int)
	for i := 0; i < GoFlushNum; i++ {
		go func(channel chan int) {
			for id := range channel {
				item := blockBuffer.blockMap[id]
				item.Lock()
				item.flush()
				item.reset()
				item.Unlock()
			}
		}(flushChannel)
	}
	for index, item := range blockBuffer.blockMap {
		if item.isDirty {
			flushChannel <- index
		}
	}
	return true, nil
}

// BeginBlockFlush 每次结束一条指令， channel接受指令并且开始刷新
func BeginBlockFlush(channel chan struct{}) {
	for _ = range channel {
		_, err := BlockFlushAll()
		if err != nil {
			fmt.Println(err)
		}
	}
}

// DeleteOldBlock
func DeleteOldBlock(filename string) error {
	blockBuffer.Lock()
	defer blockBuffer.Unlock()
	for index, item := range blockBuffer.blockMap {
		if item.filename == filename {
			item.Lock()
			delete(blockBuffer.blockMap, index)
			blockBuffer.root.remove(item)
			item.Unlock()
		}
	}
	return nil
}

// 将filename和pos转换为buffer内部的int，如果不存在则创建一个int
func Query2Int(pos nameAndPos) int {
	query2IntLock.Lock()
	defer query2IntLock.Unlock()
	if index, err := fileNamePos2Int[pos]; !err {
		return index
	}
	posNum += 1
	fileNamePos2Int[pos] = posNum
	return posNum
}
