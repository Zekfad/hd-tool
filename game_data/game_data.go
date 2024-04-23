package game_data

type ArchiveVersion uint32

const (
	ArchiveVersionHD1 ArchiveVersion = 0xf0000004
	ArchiveVersionHD2 ArchiveVersion = 0xf0000011
)

type TypeHash uint64
type NameHash = uint64

const (
	Type_ah_bin                      TypeHash = 0x2A0A70ACFE476E1D // HD2 only
	Type_animation                   TypeHash = 0x931E336D7646CC26
	Type_bik                         TypeHash = 0xAA5965F03029FA18
	Type_bones                       TypeHash = 0x18DEAD01056B72E9
	Type_camera_shake                TypeHash = 0xFCAAF813B4D3CC1E // HD2 only
	Type_cloth                       TypeHash = 0xD7014A50477953E0 // HD2 only
	Type_config                      TypeHash = 0x82645835E6B73232
	Type_flow                        TypeHash = 0x92D3EE038EEB610D // HD1 only
	Type_entity                      TypeHash = 0x9831CA893B0D087D // HD2 only
	Type_font                        TypeHash = 0x9EFE0A916AAE7880
	Type_geleta                      TypeHash = 0xB8FD4D2CEDE20ED7 // HD2 only
	Type_geometry_group              TypeHash = 0xC4F0F4BE7FB0C8D6 // HD2 only
	Type_hash_lookup                 TypeHash = 0xE3F2851035957AF5 // HD2 only
	Type_havok_ai_properties         TypeHash = 0x6592B918E67F082C // HD2 only
	Type_havok_physics_properties    TypeHash = 0xF7A09F8BB35A1D49 // HD2 only
	Type_ik_skeleton                 TypeHash = 0x57A13425279979D7 // HD2 only
	Type_level                       TypeHash = 0x2A690FD348FE9AC5
	Type_lua                         TypeHash = 0xA14E8DFA2CD117E2
	Type_material                    TypeHash = 0xEAC0B497876ADEDF
	Type_mouse_cursor                TypeHash = 0xB277B11FE4A61D37
	Type_network_config              TypeHash = 0x3B1FA9E8F6BAC374
	Type_package                     TypeHash = 0xAD9C6D9ED1E5E77A
	Type_particles                   TypeHash = 0xA8193123526FAD64
	Type_physics                     TypeHash = 0x5F7203C8F280DAB8 // HD2 only
	Type_physics_properties          TypeHash = 0xBF21403A3AB0BBB1 // HD1 only
	Type_prefab                      TypeHash = 0xAB2F78E885F513C6 // HD2 only
	Type_ragdoll_profile             TypeHash = 0x1D59BD6687DB6B33 // HD2 only
	Type_render_config               TypeHash = 0x27862FE24795319C
	Type_renderable                  TypeHash = 0x7910103158FC1DE9 // HD2 only
	Type_runtime_font                TypeHash = 0x05106B81DCD58A13 // HD2 only
	Type_shader_library              TypeHash = 0xE5EE32A477239A93
	Type_shader_library_group        TypeHash = 0x9E5C3CC74575AEB5
	Type_shading_environment         TypeHash = 0xFE73C7DCFF8A7CA5
	Type_shading_environment_mapping TypeHash = 0x250E0A11AC8E26F8 // HD2 only
	Type_speedtree                   TypeHash = 0xE985C5F61C169997 // HD2 only
	Type_state_machine               TypeHash = 0xA486D4045106165C
	Type_strings                     TypeHash = 0x0D972BAB10B40FD3
	Type_surface_properties          TypeHash = 0xAD2D3FA30D9AB394 // HD1 only
	Type_texture                     TypeHash = 0xCD4238C6A0C69E32
	Type_texture_atlas               TypeHash = 0x9199BB50B6896F02 // HD2 only
	Type_timpani_bank                TypeHash = 0x99736BE1FFF739A4 // HD1 only
	Type_timpani_master              TypeHash = 0x00A3E6C59A2B9C6C // HD1 only
	Type_unit                        TypeHash = 0xE0A48D0BE9A7453F
	Type_vector_field                TypeHash = 0xF7505933166D6755 // HD2 only
	Type_wwise_bank                  TypeHash = 0x535A7BD3E650D799 // HD2 only
	Type_wwise_dep                   TypeHash = 0xAF32095C82F2B070 // HD2 only
	Type_wwise_metadata              TypeHash = 0xD50A8B7E1C82B110 // HD2 only
	Type_wwise_properties            TypeHash = 0x5FDD5FE391076F9F // HD2 only
	Type_wwise_stream                TypeHash = 0x504B55235D21440E // HD2 only
)

type Type interface {
	GetType() TypeHash
}

type File interface {
	GetName() NameHash
	GetType() TypeHash
	GetInlineBuffer() []byte
}

type Archive interface {
	GetVersion() ArchiveVersion
	GetChecksum() uint32
	GetTypes() []Type
	GetFiles() []File
}
