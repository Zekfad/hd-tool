package game_data

import (
	"bytes"
	"encoding/binary"

	"github.com/ghostiam/binstruct"
)

type LuaFormat = uint32

const (
	LuaFormatSource          LuaFormat = 0x0
	LuaFormatGenericBytecode LuaFormat = 0x1
	LuaFormatLuaJIT2         LuaFormat = 0x2
	LuaFormatBadFormat       LuaFormat = 0x3
)

type LuaResource struct {
	Size   uint32
	Format LuaFormat
	Data   []byte `bin:"len:Size"`
}

func LuaResourceFromBytes(data []byte) (*LuaResource, error) {
	reader := binstruct.NewReaderFromBytes(data, binary.LittleEndian, false)
	var resource LuaResource
	err := reader.Unmarshal(&resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (resource LuaResource) ToBytes() ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, uint32(len(resource.Data)))
	binary.Write(b, binary.LittleEndian, uint32(resource.Format))
	b.Write(resource.Data)
	return b.Bytes(), nil
}
