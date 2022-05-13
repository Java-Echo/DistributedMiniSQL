package Buffermanager

import (
	"io"
	"os"
	"sync"
)

type block struct {
	filename string
	blockid  uint16
	isDirty  bool
	pin      bool
	Data     []byte
	next     *block
	prev     *block // double linked list
	sync.Mutex
}

// operations
// set dirty page
func (b *block) setDirty() {
	b.isDirty = true
}

// set pin block
func (b *block) setPin() {
	b.pin = true
}

// unset pin block
func (b *block) setUnPin() {
	b.pin = false
}

// release the lock of reading block
func (b *block) finishRead() {
	b.Unlock()
	return
}

//reset the block
func (b *block) reset() {
	b.isDirty = false
	b.pin = false
}

// init the block
func (b *block) init(filename string, bid uint16) {
	b.filename = filename
	b.blockid = bid
}

// read the file
func (b *block) read() error {
	if b.isDirty {
		return b.flush()
	}
	file, err := os.Open(b.filename)
	if err != nil {
		return nil
	}
	defer file.Close()
	if err != nil {
		return err
	}
	bid64 := int64(b.blockid)
	_, err = file.Seek(bid64*blockSize, 0)
	if err != nil {
		return err
	}

	_, err = io.ReadFull(file, b.Data)
	if err != nil {
		return err
	}

	return err
}

// read back and flush
func (b *block) flush() error {
	if !b.isDirty {
		return nil
	}
	file, err := os.OpenFile(b.filename, os.O_WRONLY, 0666) // read and write
	defer file.Close()
	if err != nil {
		return err
	}
	bid64 := int64(b.blockid)
	_, err = file.Seek(bid64*blockSize, 0) // _ is a placeholder
	if err != nil {
		return err
	}
	_, err = file.Write(b.Data)
	b.isDirty = false
	return err
}
