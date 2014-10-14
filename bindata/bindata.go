package bindata

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

func res_sqlite_deltaiota_db() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0xec, 0x97,
		0xcf, 0x4e, 0xdb, 0x4e, 0x10, 0xc7, 0xbd, 0x76, 0x82, 0x21, 0xe8, 0xf7,
		0xbb, 0xf4, 0xe0, 0x13, 0xd2, 0xca, 0x17, 0x62, 0x01, 0x07, 0xca, 0xa1,
		0xd7, 0xa6, 0x74, 0x41, 0x11, 0xc1, 0x01, 0xe3, 0x48, 0x70, 0x8a, 0x5c,
		0x58, 0xda, 0xa5, 0xb1, 0x1d, 0xbc, 0x8e, 0x9a, 0x1e, 0x43, 0x79, 0x8e,
		0x3e, 0x49, 0x9f, 0xa3, 0x8f, 0xd0, 0x97, 0xe8, 0xa5, 0xde, 0x75, 0x1c,
		0xe2, 0xc4, 0x21, 0xca, 0x29, 0x07, 0xcf, 0x47, 0xf2, 0x9f, 0x9d, 0x99,
		0xf5, 0x8c, 0x77, 0xf6, 0xeb, 0x28, 0x57, 0x97, 0x2d, 0x16, 0x53, 0x7c,
		0x1f, 0x46, 0xbe, 0x17, 0xe3, 0x23, 0xa5, 0xa2, 0x20, 0xa4, 0xbc, 0xc7,
		0x58, 0x51, 0x14, 0x94, 0x1c, 0xdb, 0xca, 0x0b, 0x7a, 0x72, 0x54, 0xa6,
		0xc6, 0x48, 0x59, 0x0e, 0x52, 0x0e, 0xfe, 0xa8, 0x55, 0x71, 0xa3, 0xfd,
		0x15, 0xe3, 0xcd, 0xf4, 0x02, 0x00, 0xc0, 0xfa, 0xd1, 0xb5, 0xff, 0xc4,
		0xa5, 0xb2, 0xee, 0x3a, 0x00, 0x00, 0x58, 0x07, 0xa0, 0x7f, 0x00, 0x28,
		0x2f, 0xa0, 0x7f, 0x00, 0x28, 0x2f, 0x35, 0x71, 0x02, 0xfd, 0x03, 0x40,
		0x29, 0x81, 0xdf, 0x7f, 0x00, 0x28, 0x2f, 0x42, 0xff, 0x1a, 0xe2, 0x0a,
		0xe2, 0xea, 0x2f, 0xed, 0xff, 0x75, 0x57, 0x53, 0x0a, 0x9e, 0x4e, 0x90,
		0x6e, 0xec, 0xee, 0xa2, 0xe7, 0xb3, 0xd8, 0xfb, 0xd4, 0xa3, 0x41, 0x18,
		0xb3, 0x7b, 0x76, 0xeb, 0xc5, 0x2c, 0x0c, 0x78, 0x6e, 0xa0, 0x1e, 0x3b,
		0xa4, 0xe1, 0x12, 0xec, 0x36, 0x3e, 0xb4, 0x08, 0x36, 0x73, 0x3e, 0x13,
		0xd7, 0x6b, 0x5b, 0x26, 0xbb, 0x33, 0xf1, 0x84, 0xa6, 0xed, 0x92, 0x53,
		0xe2, 0xe0, 0x0b, 0xa7, 0x79, 0xde, 0x70, 0x6e, 0xf0, 0x19, 0xb9, 0xc1,
		0x8d, 0x8e, 0xdb, 0x6e, 0xda, 0xc9, 0x73, 0xce, 0x89, 0xed, 0xd6, 0xb6,
		0xf6, 0xb1, 0x39, 0xe0, 0x34, 0xea, 0xa6, 0xf3, 0xb2, 0x09, 0x76, 0xdb,
		0xc5, 0x76, 0xa7, 0xd5, 0x92, 0xfe, 0x98, 0xf9, 0x94, 0xc7, 0x9e, 0xdf,
		0x37, 0x8b, 0xfd, 0x11, 0xf5, 0xb2, 0xa4, 0xc5, 0xf3, 0xe9, 0x30, 0x7e,
		0x29, 0xca, 0x25, 0xd7, 0x6e, 0xde, 0x3f, 0x88, 0xd8, 0x54, 0xcd, 0x79,
		0xbf, 0x08, 0x38, 0x69, 0x3b, 0xa4, 0x79, 0x6a, 0x8b, 0xea, 0xeb, 0xe3,
		0x5a, 0x2d, 0xec, 0x90, 0x13, 0xe2, 0x10, 0xfb, 0x98, 0x5c, 0x61, 0x61,
		0xe3, 0xf5, 0xc4, 0x58, 0xb3, 0x2e, 0xd4, 0x0d, 0x63, 0x6f, 0x0f, 0xdd,
		0xc8, 0x65, 0xe4, 0x8f, 0x3d, 0x16, 0xd3, 0x2e, 0xa7, 0x8f, 0x03, 0x1a,
		0xdc, 0xce, 0x0e, 0xb5, 0xdc, 0x52, 0xce, 0x38, 0xeb, 0x81, 0xe7, 0xd3,
		0xfd, 0x64, 0x64, 0x8d, 0x42, 0x4d, 0x37, 0x76, 0x76, 0xd0, 0x8f, 0xc3,
		0xf4, 0x99, 0x94, 0x73, 0xb1, 0xd8, 0xd9, 0xb5, 0x92, 0x6f, 0x48, 0x66,
		0x9e, 0xeb, 0xc5, 0x8a, 0x9d, 0x28, 0x5c, 0xc7, 0xaf, 0xf4, 0xbb, 0xb9,
		0x78, 0x15, 0xe9, 0xb0, 0xcf, 0x22, 0x6a, 0x16, 0xf4, 0x60, 0xc5, 0x35,
		0x14, 0xfa, 0xaf, 0xa2, 0xb7, 0x8a, 0xea, 0xab, 0xbf, 0x55, 0x84, 0x7e,
		0x26, 0xb7, 0xc0, 0xab, 0x3c, 0x6c, 0xea, 0xc6, 0x3b, 0x03, 0x8d, 0xde,
		0xb0, 0xe0, 0x8e, 0x0e, 0xe5, 0x52, 0x76, 0x07, 0x01, 0x4b, 0xb6, 0x52,
		0xb7, 0xef, 0x71, 0xfe, 0x2d, 0x8c, 0xee, 0xa4, 0x71, 0x7b, 0xbc, 0x59,
		0x3a, 0x76, 0xf3, 0xb2, 0x43, 0x92, 0x46, 0x7d, 0x24, 0xd7, 0x69, 0xd3,
		0xe7, 0xe2, 0x4d, 0xdc, 0xb6, 0xc7, 0xae, 0x64, 0x2f, 0x99, 0x13, 0xb3,
		0xe5, 0xe9, 0xba, 0x71, 0x98, 0xe4, 0xd2, 0xe7, 0x73, 0x51, 0xdf, 0x63,
		0x3d, 0x69, 0xa9, 0x2d, 0x4f, 0x24, 0x83, 0x67, 0xb2, 0xa4, 0x36, 0xeb,
		0x61, 0x63, 0xe1, 0xeb, 0x88, 0x81, 0x90, 0x86, 0x34, 0x6e, 0x2d, 0xcf,
		0x92, 0xc5, 0xcf, 0x24, 0x9a, 0x98, 0xad, 0xcf, 0x15, 0xdd, 0x38, 0xda,
		0x41, 0xa3, 0x6d, 0x99, 0x2b, 0x93, 0x4f, 0x36, 0x3d, 0xd9, 0xf1, 0x99,
		0xa9, 0x5a, 0x98, 0xac, 0x60, 0x42, 0x9a, 0x6a, 0x4a, 0x88, 0x52, 0x37,
		0xd6, 0xd3, 0x6e, 0x55, 0x37, 0x0c, 0x03, 0x3d, 0x1f, 0x48, 0x1d, 0xcb,
		0x52, 0xe4, 0x69, 0x23, 0xaf, 0xe0, 0xac, 0xc6, 0xd9, 0x4f, 0xe9, 0x0a,
		0x0a, 0x4e, 0x5f, 0x78, 0x81, 0x4c, 0xef, 0x59, 0xc4, 0xe3, 0xee, 0x24,
		0x64, 0x3e, 0xa0, 0xe7, 0x4d, 0xfb, 0x8b, 0x84, 0x9e, 0x76, 0x6e, 0xf1,
		0x97, 0xa0, 0xff, 0x25, 0x0c, 0xe8, 0xab, 0x01, 0x93, 0x4d, 0x56, 0x10,
		0x60, 0xc1, 0xff, 0x7f, 0x00, 0x28, 0x2f, 0xa0, 0x7f, 0x00, 0x28, 0x2f,
		0xa0, 0x7f, 0x00, 0x28, 0x2f, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x73,
		0x25, 0x27, 0x05, 0x00, 0x2c, 0x00, 0x00,
	},
		"res/sqlite/deltaiota.db",
	)
}

func res_sqlite_deltaiota_sql() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x94, 0x92,
		0xcf, 0xce, 0xa2, 0x30, 0x14, 0xc5, 0xd7, 0x1f, 0x4f, 0x71, 0xd3, 0x95,
		0x9a, 0x49, 0x7c, 0x00, 0x57, 0xe8, 0x5c, 0x4d, 0x33, 0x52, 0x66, 0x6a,
		0x49, 0x74, 0x45, 0x1a, 0xad, 0xb1, 0x19, 0xfe, 0x28, 0x2d, 0x19, 0xe7,
		0xed, 0x3f, 0x40, 0x11, 0x89, 0xd5, 0x28, 0xcb, 0x7b, 0xee, 0x3d, 0x9c,
		0xf3, 0x4b, 0xc7, 0x23, 0xd8, 0xa9, 0xc4, 0x4a, 0x9d, 0x5b, 0x09, 0xe6,
		0x94, 0x68, 0xab, 0xc0, 0x6c, 0x0f, 0x2a, 0x95, 0x30, 0x1a, 0x7b, 0x53,
		0x5c, 0x50, 0x06, 0x82, 0xfb, 0x6c, 0xe5, 0xcf, 0x04, 0x0d, 0xd9, 0xc4,
		0xab, 0x0e, 0xb2, 0xdc, 0xea, 0xbd, 0xde, 0x4a, 0xab, 0xf3, 0xcc, 0xd4,
		0x6b, 0x33, 0x8e, 0xbe, 0x40, 0x10, 0xfe, 0x74, 0x89, 0x40, 0x7a, 0x32,
		0x81, 0x81, 0xf7, 0x45, 0xf4, 0x8e, 0xc0, 0xed, 0xa3, 0x4c, 0xe0, 0x02,
		0x39, 0xfc, 0xe6, 0x34, 0xf0, 0xf9, 0x06, 0x7e, 0xe1, 0x06, 0xfc, 0x48,
		0x84, 0x94, 0x55, 0x3e, 0x01, 0x32, 0xe1, 0x7d, 0xfd, 0x00, 0x52, 0x1a,
		0x55, 0xc4, 0x97, 0xbb, 0xf6, 0x80, 0x85, 0x02, 0x58, 0xb4, 0x5c, 0x36,
		0xba, 0xd5, 0xa9, 0x32, 0x56, 0xa6, 0x47, 0xe2, 0xd6, 0x0b, 0x25, 0xdb,
		0x9f, 0xba, 0xef, 0xd5, 0xd9, 0x76, 0xa1, 0x04, 0xae, 0x45, 0x5f, 0x2f,
		0x0b, 0x7d, 0x97, 0xb9, 0xaf, 0xd7, 0x0b, 0xf3, 0x90, 0x23, 0x5d, 0xb0,
		0x3a, 0xfd, 0xe0, 0x9a, 0x75, 0x08, 0x1c, 0xe7, 0xc8, 0x91, 0xcd, 0x70,
		0x05, 0xf5, 0xcc, 0x0c, 0xaa, 0xa1, 0x37, 0x6c, 0xa0, 0x19, 0x65, 0x8c,
		0x9b, 0x57, 0xab, 0x3c, 0xa0, 0xfa, 0x10, 0x94, 0xb3, 0xe6, 0x5f, 0xf5,
		0x9f, 0x3c, 0x2f, 0xa9, 0xce, 0x47, 0x5d, 0x28, 0xe2, 0x40, 0xf4, 0x71,
		0xc5, 0x6b, 0xa7, 0x88, 0xd1, 0x3f, 0x11, 0x56, 0x7e, 0x3f, 0x71, 0xdd,
		0x55, 0x8b, 0xcb, 0x4c, 0x9f, 0x4a, 0x15, 0x37, 0x69, 0x42, 0xd6, 0xeb,
		0xdc, 0x44, 0xbc, 0x30, 0x6a, 0x0c, 0x1f, 0x01, 0x35, 0x63, 0xc7, 0x43,
		0xfa, 0x00, 0x50, 0x26, 0x53, 0x45, 0x9e, 0x51, 0xd8, 0xeb, 0xc2, 0xd8,
		0xf8, 0xb6, 0xf2, 0xb8, 0x90, 0xc8, 0x7b, 0xdd, 0xc5, 0x31, 0x95, 0x3a,
		0x79, 0xf5, 0x9a, 0x8e, 0x87, 0x3c, 0x53, 0x2f, 0x17, 0xa4, 0x31, 0xff,
		0xf2, 0x62, 0xe7, 0x0c, 0xf9, 0x0c, 0x6f, 0x03, 0xa6, 0x65, 0xdb, 0xd5,
		0xac, 0x01, 0xb7, 0xcc, 0xba, 0xf6, 0x6f, 0x99, 0x5c, 0x8b, 0xf4, 0x1c,
		0x2e, 0xb3, 0xb7, 0xce, 0xbb, 0x16, 0x3d, 0x87, 0xdb, 0xb8, 0x36, 0x09,
		0x83, 0x80, 0x8a, 0x89, 0xf7, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x71, 0x08,
		0x13, 0xe1, 0x79, 0x04, 0x00, 0x00,
	},
		"res/sqlite/deltaiota.sql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"res/sqlite/deltaiota.db": res_sqlite_deltaiota_db,
	"res/sqlite/deltaiota.sql": res_sqlite_deltaiota_sql,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"res": &_bintree_t{nil, map[string]*_bintree_t{
		"sqlite": &_bintree_t{nil, map[string]*_bintree_t{
			"deltaiota.db": &_bintree_t{res_sqlite_deltaiota_db, map[string]*_bintree_t{
			}},
			"deltaiota.sql": &_bintree_t{res_sqlite_deltaiota_sql, map[string]*_bintree_t{
			}},
		}},
	}},
}}
