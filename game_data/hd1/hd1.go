package hd1

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"

	"io"

	"github.com/Zekfad/hd-tool/game_data"
	"github.com/ghostiam/binstruct"
)

/**
 * Packed container
 */

const CompressedChunkSize = 65536

type Archive struct {
	ArchiveVersion game_data.ArchiveVersion `bin:"ReadVersion"`
	UnpackedSize   uint32
	Reserved       uint32        // must be 0
	Chunks         []PackedChunk `bin:"ReadChunks"`

	Unpacked UnpackedArchive `bin:"UnpackArchive"`
}

type PackedChunk struct {
	Size uint32
	Data []byte `bin:"len:Size"`
}

func (chunk PackedChunk) IsCompressed() bool {
	return chunk.Size != CompressedChunkSize
}

func (data *Archive) ReadChunks(r binstruct.Reader) error {
	data.Chunks = make([]PackedChunk, 0)

	var chunk PackedChunk
	for {
		err := r.Unmarshal(&chunk)
		eof := errors.Is(err, io.EOF)
		if err != nil && !eof {
			return err
		}
		data.Chunks = append(data.Chunks, chunk)
		if eof {
			break
		}
	}
	return nil
}

func (archive *Archive) UnpackArchive(r binstruct.Reader) error {
	buffer := make([]byte, 0, archive.UnpackedSize)
	for _, packed := range archive.Chunks {
		if packed.IsCompressed() {
			z, err := zlib.NewReader(bytes.NewReader(packed.Data))
			if err != nil {
				return err
			}
			defer z.Close()
			data, err := io.ReadAll(z)
			if err != nil {
				return err
			}
			buffer = append(buffer, data...)
		} else {
			buffer = append(buffer, packed.Data...)
		}
	}
	reader := binstruct.NewReaderFromBytes(buffer, binary.LittleEndian, false)
	return reader.Unmarshal(&archive.Unpacked)
}

/**
 * Exploded package
 */

type UnpackedArchive struct {
	Header ArchiveHeader
	Types  []Type `bin:"len:Header.EntriesCount"`
	Files  []File `bin:"len:Header.EntriesCount"`
}

type ArchiveHeader struct {
	EntriesCount uint32
	Magic        []byte `bin:"len:256"`
}

type Type struct {
	Type game_data.TypeHash
	Name game_data.NameHash
}

type VariantHeader struct {
	Unk00      uint32
	Size       uint32
	StreamSize uint32
}

type File struct {
	Type game_data.TypeHash
	Name game_data.NameHash

	VariantsCount uint32
	StreamOffset  uint32

	VariantHeaders []VariantHeader `bin:"len:VariantsCount"`
	VariantBuffers [][]byte        `bin:"ReadBuffers"`
}

func (header *Archive) ReadVersion(r binstruct.Reader) error {
	_version, err := r.ReadUint32()
	if err != nil {
		return err
	}
	switch version := game_data.ArchiveVersion(_version); version {
	case game_data.ArchiveVersionHD1:
		header.ArchiveVersion = version
		return nil
	default:
		return fmt.Errorf("invalid archive version: %#08X", version)
	}
}

func (file *File) ReadBuffers(r binstruct.Reader) error {
	file.VariantBuffers = make([][]byte, file.VariantsCount)
	for i, variant := range file.VariantHeaders {
		size := int(variant.Size)

		read, buffer, err := r.ReadBytes(size)
		if err != nil {
			return err
		}
		if read != size {
			return fmt.Errorf("not enough data to read variant buffer: read %d, expected %d", read, size)
		}
		file.VariantBuffers[i] = buffer
	}
	return nil
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
	return bytes.Join(file.VariantBuffers, nil)
}

/**
 * Archive interface
 */

// GetVersion implements Archive.
func (archive Archive) GetVersion() game_data.ArchiveVersion {
	return archive.ArchiveVersion
}

// GetChecksum implements Archive.
func (archive Archive) GetChecksum() uint32 {
	return 0
}

// GetTypes implements Archive.
func (archive Archive) GetTypes() []game_data.Type {
	types := make([]game_data.Type, len(archive.Unpacked.Types))
	for i, d := range archive.Unpacked.Types {
		types[i] = d
	}
	return types
}

// GetFiles implements Archive.
func (archive Archive) GetFiles() []game_data.File {
	files := make([]game_data.File, len(archive.Unpacked.Files))
	for i, d := range archive.Unpacked.Files {
		files[i] = d
	}
	return files
}
