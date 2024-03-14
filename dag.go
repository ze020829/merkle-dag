package merkledag

import "hash"

// 链接结构体
type Link struct {
	Name string
	Hash []byte
	Size int
}

// 对象结构体
type Object struct {
	Links []Link
	Data  []byte
}

func Add(store KVStore, node Node, h hash.Hash) ([]byte, error) {
	// 检查节点是否为文件类型
	file, ok := node.(File)
	if !ok {
		return nil, nil
	}

	// 获取文件的字节数据
	data := file.Bytes()

	// 分片大小
	chunkSize := 256

	// 存储所有分片的哈希值
	var hashes [][]byte

	// 将文件数据分片，并计算每个分片的哈希值
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize

		// 如果结束索引超过了数据的长度，
		if end > len(data) {
			end = len(data)
		}

		// 获取当前分片
		chunk := data[i:end]

		// 计算分片的哈希值
		h.Reset()
		h.Write(chunk)
		hash := h.Sum(nil)

		// 将哈希值存储到KVStore中
		key := hash
		store.Put(key, chunk)

		// 将哈希值添加到哈希列表中
		hashes = append(hashes, hash)
	}

	// 创建一个新的对象
	object := &Object{
		Data: data,
	}

	// 将所有分片的哈希值添加到对象的链接中
	for _, hash := range hashes {
		link := Link{
			Name: "",
			Hash: hash,
			Size: len(hash),
		}
		object.Links = append(object.Links, link)
	}

	// 计算对象的Merkle Root
	h.Reset()
	h.Write(object.Data)
	merkleRoot := h.Sum(nil)

	// 返回Merkle Root
	return merkleRoot, nil
}
