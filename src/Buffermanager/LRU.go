package Buffermanager

import "sync"

const InitSize = 1024
const DeleteSize = 256 // when the Cache is full

type LRUList struct {
	root block // using dummy header
	len  int
}

func NewLRUList() *LRUList {
	list := new(LRUList)
	// list is a loop
	list.root.next = &list.root
	list.root.prev = &list.root
	return list
}

//return head

func (l *LRUList) Front() *block {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// insert a after b

func (l *LRUList) insert(a, b *block) *block {
	tmp := b.next
	b.next = a
	a.prev = b
	a.next = tmp
	tmp.prev = a
	l.len += 1
	return a
}

// remove a

func (l *LRUList) remove(a *block) *block {
	a.prev.next = a.next
	a.next.prev = a.prev
	l.len -= 1
	return a // a is collected as garbage by GC (c++: ?
}

// move the recent accessed block to back
func (l *LRUList) moveToBack(a *block) {
	if l.root.prev == a { // if it has only 1 block
		return
	}
	l.insert(l.remove(a), l.root.prev)
}

// append a new block
func (l *LRUList) append(a *block) {
	l.insert(a, l.root.prev)
}

//using list and map to achieve LRU

type LRUCache struct {
	Capacity int
	root     *LRUList
	blockMap map[int]*block
	sync.Mutex
}

// new

func NewLRUCache() *LRUCache {
	cache := new(LRUCache)
	cache.Capacity = InitSize
	cache.root = NewLRUList()
	cache.blockMap = make(map[int]*block, InitSize*2)
	return cache
}

// put block int buffer
func (cache *LRUCache) insertBlock(val *block, index int) *block {
	cache.Lock()
	defer cache.Unlock()
	if item, err := cache.blockMap[index]; !err {
		cache.root.moveToBack(item)
		return item
	}
	if len(cache.blockMap) >= cache.Capacity {
		var tmp = cache.root.Front()
		for _i := 0; _i < DeleteSize; _i++ {
			if tmp.pin {
				tmp = tmp.next
			} else {
				tmp.Lock()
				tmp.flush()
				cache.root.remove(tmp)
				oldIndex := Query2Int(nameAndPos{filename: tmp.filename, blockid: tmp.blockid})
				delete(cache.blockMap, oldIndex)
				tmp.Unlock()
				tmp = tmp.next
			}
		}
	}
	cache.root.append(val)
	cache.blockMap[index] = val
	return val
}

// map index to block
func (cache *LRUCache) findBlock(index int) (bool, *block) {
	cache.Lock()
	defer cache.Unlock()
	if ret, err := cache.blockMap[index]; !err {
		cache.root.moveToBack(ret)
		return true, ret
	}
	return false, nil
}
