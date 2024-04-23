package writer_hd1

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"

	"github.com/Zekfad/hd-tool/game_data/hd1"
)

var le = binary.LittleEndian

func WriteArchive(archive hd1.Archive, writer io.Writer) error {
	binary.Write(writer, le, uint32(archive.ArchiveVersion))

	buffer := new(bytes.Buffer)
	binary.Write(buffer, le, archive.Unpacked.Header.EntriesCount)
	buffer.Write(archive.Unpacked.Header.Magic)
	for _, _type := range archive.Unpacked.Types {
		binary.Write(buffer, le, uint64(_type.Type))
		binary.Write(buffer, le, uint64(_type.Name))
	}
	for _, file := range archive.Unpacked.Files {
		binary.Write(buffer, le, uint64(file.Type))
		binary.Write(buffer, le, uint64(file.Name))
		binary.Write(buffer, le, file.VariantsCount)
		binary.Write(buffer, le, file.StreamOffset)

		for i, header := range file.VariantHeaders {
			binary.Write(buffer, le, header.Unk00)
			// binary.Write(buffer, le, header.Size)
			binary.Write(buffer, le, uint32(len(file.VariantBuffers[i])))
			binary.Write(buffer, le, header.StreamSize)
		}
		for _, variantBuffer := range file.VariantBuffers {
			buffer.Write(variantBuffer)
		}
	}

	dataRaw := buffer.Bytes()

	binary.Write(writer, le, uint32(len(dataRaw)))
	binary.Write(writer, le, uint32(0))

	for i := 0; i < len(dataRaw); i += hd1.CompressedChunkSize {
		part := dataRaw[i : i+hd1.CompressedChunkSize]
		b := new(bytes.Buffer)
		z := zlib.NewWriter(b)
		z.Write(part)
		z.Close()

		compressed := b.Bytes()

		binary.Write(writer, le, uint32(len(compressed)))
		writer.Write(compressed)
	}
	return nil
}
