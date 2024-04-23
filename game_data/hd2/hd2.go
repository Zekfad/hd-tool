package hd2

import (
	"fmt"

	"github.com/Zekfad/hd-tool/game_data"
	"github.com/ghostiam/binstruct"
)

type Archive struct {
	Header ArchiveHeader
	Types  []Type `bin:"len:Header.TypesCount"`
	Files  []File `bin:"len:Header.FilesCount"`
}

type ArchiveHeader struct {
	ArchiveVersion game_data.ArchiveVersion `bin:"ReadVersion"`
	TypesCount     uint32
	FilesCount     uint32

	Unk00         uint32
	Checksum      uint32
	Unk01         uint32
	Unk02         uint64 // looks like a hash, but isn't
	BufferSize    uint64
	GpuBufferSize uint64

	Unk03 uint32
	Unk04 uint32
	Unk05 uint32
	Unk06 uint32
	Unk07 uint32
	Unk08 uint32
}

type Type struct {
	Unk00        uint32
	Unk01        uint32
	Type         game_data.TypeHash
	Count        uint32
	Unk02        uint32
	Alignment    uint32
	GpuAlignment uint32
}

type File struct {
	Name game_data.NameHash
	Type game_data.TypeHash

	// aligned to alignment
	Offset uint64
	// unaligned raw aux buffer
	StreamOffset uint64
	// aligned to gpu_alignment
	GpuOffset uint64

	// aligned to 0x200
	// seems to be an offset of read buffer that have size of 0x200 per element
	BufferOffset uint64
	// aligned to 0x600
	// seems to be an offset of read buffer that have size of 0x600 per element
	GpuBufferOffset uint64

	Size          uint32
	StreamSize    uint32
	GpuStreamSize uint32

	Alignment    uint32
	GpuAlignment uint32
	Index        uint32

	InlineBuffer []byte `bin:"offsetStart:Offset, len:Size, offsetRestore"`
}

func (header *ArchiveHeader) ReadVersion(r binstruct.Reader) error {
	_version, err := r.ReadUint32()
	if err != nil {
		return err
	}
	switch version := game_data.ArchiveVersion(_version); version {
	case game_data.ArchiveVersionHD2:
		header.ArchiveVersion = version
		return nil
	default:
		return fmt.Errorf("invalid archive version: %#08X", version)
	}
}

/**
 * Type interface
 */

// GetType implements Type.
func (_type Type) GetType() game_data.TypeHash {
	return _type.Type
}

/**
 * File interface
 */

// GetName implements File.
func (file File) GetName() game_data.NameHash {
	return file.Name
}

// GetType implements File.
func (file File) GetType() game_data.TypeHash {
	return file.Type
}

// GetInlineBuffer implements File.
func (file File) GetInlineBuffer() []byte {
	return file.InlineBuffer
}

/**
 * Archive interface
 */

// GetVersion implements Archive.
func (archive Archive) GetVersion() game_data.ArchiveVersion {
	return archive.Header.ArchiveVersion
}

// GetChecksum implements Archive.
func (archive Archive) GetChecksum() uint32 {
	return archive.Header.Checksum
}

// GetTypes implements Archive.
func (archive Archive) GetTypes() []game_data.Type {
	types := make([]game_data.Type, len(archive.Types))
	for i, d := range archive.Types {
		types[i] = d
	}
	return types
}

// GetFiles implements Archive.
func (archive Archive) GetFiles() []game_data.File {
	files := make([]game_data.File, len(archive.Files))
	for i, d := range archive.Files {
		files[i] = d
	}
	return files
}
