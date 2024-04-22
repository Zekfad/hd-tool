package game_data

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"

	"github.com/ghostiam/binstruct"
)

type ArchiveVersion uint32

const (
	ArchiveVersionHD1 ArchiveVersion = 0xf0000004
	ArchiveVersionHD2 ArchiveVersion = 0xf0000011
)

type TypeHash uint64
type NameHash = uint64

const (
	Type_ah_bin                      TypeHash = 0x2A0A70ACFE476E1D
	Type_animation                   TypeHash = 0x931E336D7646CC26
	Type_bik                         TypeHash = 0xAA5965F03029FA18
	Type_bones                       TypeHash = 0x18DEAD01056B72E9
	Type_camera_shake                TypeHash = 0xFCAAF813B4D3CC1E
	Type_cloth                       TypeHash = 0xD7014A50477953E0
	Type_config                      TypeHash = 0x82645835E6B73232
	Type_entity                      TypeHash = 0x9831CA893B0D087D
	Type_font                        TypeHash = 0x9EFE0A916AAE7880
	Type_geleta                      TypeHash = 0xB8FD4D2CEDE20ED7
	Type_geometry_group              TypeHash = 0xC4F0F4BE7FB0C8D6
	Type_hash_lookup                 TypeHash = 0xE3F2851035957AF5
	Type_havok_ai_properties         TypeHash = 0x6592B918E67F082C
	Type_havok_physics_properties    TypeHash = 0xF7A09F8BB35A1D49
	Type_ik_skeleton                 TypeHash = 0x57A13425279979D7
	Type_level                       TypeHash = 0x2A690FD348FE9AC5
	Type_lua                         TypeHash = 0xA14E8DFA2CD117E2
	Type_material                    TypeHash = 0xEAC0B497876ADEDF
	Type_mouse_cursor                TypeHash = 0xB277B11FE4A61D37
	Type_network_config              TypeHash = 0x3B1FA9E8F6BAC374
	Type_package                     TypeHash = 0xAD9C6D9ED1E5E77A
	Type_particles                   TypeHash = 0xA8193123526FAD64
	Type_physics                     TypeHash = 0x5F7203C8F280DAB8
	Type_prefab                      TypeHash = 0xAB2F78E885F513C6
	Type_ragdoll_profile             TypeHash = 0x1D59BD6687DB6B33
	Type_render_config               TypeHash = 0x27862FE24795319C
	Type_renderable                  TypeHash = 0x7910103158FC1DE9
	Type_runtime_font                TypeHash = 0x05106B81DCD58A13
	Type_shader_library              TypeHash = 0xE5EE32A477239A93
	Type_shader_library_group        TypeHash = 0x9E5C3CC74575AEB5
	Type_shading_environment         TypeHash = 0xFE73C7DCFF8A7CA5
	Type_shading_environment_mapping TypeHash = 0x250E0A11AC8E26F8
	Type_speedtree                   TypeHash = 0xE985C5F61C169997
	Type_state_machine               TypeHash = 0xA486D4045106165C
	Type_strings                     TypeHash = 0x0D972BAB10B40FD3
	Type_texture                     TypeHash = 0xCD4238C6A0C69E32
	Type_texture_atlas               TypeHash = 0x9199BB50B6896F02
	Type_unit                        TypeHash = 0xE0A48D0BE9A7453F
	Type_vector_field                TypeHash = 0xF7505933166D6755
	Type_wwise_bank                  TypeHash = 0x535A7BD3E650D799
	Type_wwise_dep                   TypeHash = 0xAF32095C82F2B070
	Type_wwise_metadata              TypeHash = 0xD50A8B7E1C82B110
	Type_wwise_properties            TypeHash = 0x5FDD5FE391076F9F
	Type_wwise_stream                TypeHash = 0x504B55235D21440E
)

type Archive struct {
	Header ArchiveHeader
	Types  []Type `bin:"len:Header.TypesCount"`
	Files  []File `bin:"len:Header.FilesCount"`
}

type ArchiveHeader struct {
	ArchiveVersion ArchiveVersion `bin:"ReadVersion"`
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
	Unk02        uint64
	Count        uint32
	Unk03        uint32
	Alignment    uint32
	GpuAlignment uint32
}

type File struct {
	Name NameHash
	Type TypeHash

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
	switch version := ArchiveVersion(_version); version {
	case ArchiveVersionHD2:
		header.ArchiveVersion = version
		return nil
	default:
		return fmt.Errorf("invalid archive version: %#08X", version)
	}
}

func ArchiveFromBytes(data []byte) (*Archive, error) {
	reader := binstruct.NewReaderFromBytes(data, binary.LittleEndian, false)
	var archive Archive
	err := reader.Unmarshal(&archive)
	if err != nil {
		return nil, err
	}
	return &archive, nil
}

func ArchiveFromFile(name string) (*Archive, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to open file"),
			err,
		)
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to read file"),
			err,
		)
	}

	archive, err := ArchiveFromBytes(buffer)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to parse file"),
			err,
		)
	}
	return archive, nil
}

func ArchivesFromDirectory(dirname string) iter.Seq2[string, Archive] {
	return func(yield func(path string, archive Archive) bool) {
		entries, err := os.ReadDir(dirname)
		if err != nil {
			fmt.Printf("Failed to read directory: %s", err)
			return
		}

		for _, file := range entries {
			if file.IsDir() {
				continue
			}
			if filepath.Ext(file.Name()) == "" {
				fullPath := filepath.Join(dirname, file.Name())

				// fmt.Printf("Candidate %s ... ", fullPath)
				archive, err := ArchiveFromFile(fullPath)
				if err != nil {
					// fmt.Printf("not a package: %s\n", err)
					continue
				}
				// fmt.Printf("is a valid package!\n")
				if !yield(fullPath, *archive) {
					return
				}
			}
		}
	}
}
