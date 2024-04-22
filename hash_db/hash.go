package hash_db

import (
	"github.com/nfisher/gstream/hash/murmur2"
)

func Hash(data string) uint64 {
	return murmur2.Hash([]byte(data), 0)
}
