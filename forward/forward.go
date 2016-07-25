package forward

type WritePacker interface {
	// 注意：由于可能内部使用了队列机制
	// 即使函数安全返回，也不允许需要改 data 的值
	WritePack(data []byte) error
}

type ReadPacker interface {
	ReadPack() ([]byte, error)
}

func Forward(send WritePacker, read ReadPacker) error {
	for {
		buf, err := read.ReadPack()
		if err != nil {
			return err
		}
		if err := send.WritePack(buf); err != nil {
			return err
		}
	}
}
