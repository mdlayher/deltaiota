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
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0xec, 0x93,
		0xcf, 0x6e, 0x9b, 0x40, 0x10, 0xc6, 0x77, 0xfd, 0x0f, 0x15, 0x29, 0xa7,
		0x1c, 0xf6, 0x14, 0x69, 0xc5, 0x29, 0x28, 0xe9, 0xa1, 0xea, 0xa1, 0xea,
		0xad, 0xb4, 0x59, 0x55, 0xa8, 0x04, 0x27, 0x14, 0xa4, 0xf8, 0x84, 0x68,
		0xb3, 0x51, 0x37, 0xc5, 0xe0, 0xb0, 0xb8, 0x4d, 0x8e, 0x49, 0x1f, 0xa6,
		0xef, 0xd5, 0x43, 0x1f, 0xa0, 0x4f, 0xd0, 0x05, 0x0c, 0x8e, 0x21, 0x96,
		0xdb, 0x33, 0xf3, 0x93, 0x16, 0x7b, 0xe6, 0xdb, 0xd9, 0x61, 0x3f, 0x7b,
		0x3e, 0x9e, 0x3b, 0x22, 0xe7, 0xf4, 0x2a, 0xcd, 0xe6, 0x51, 0x4e, 0x5f,
		0xa2, 0x11, 0xc2, 0x18, 0xbd, 0xa1, 0x14, 0x21, 0x84, 0xd5, 0xd2, 0xd0,
		0x9a, 0xb1, 0x5a, 0xa3, 0x47, 0x31, 0x46, 0xbb, 0xc1, 0xe8, 0xf9, 0xef,
		0xc1, 0x9e, 0xfa, 0x32, 0x41, 0x3f, 0xd1, 0xf0, 0x64, 0xf0, 0x07, 0xff,
		0xc2, 0xdf, 0xf0, 0xbe, 0x0a, 0xfe, 0x8f, 0xeb, 0x89, 0x46, 0x5e, 0x11,
		0x7c, 0xbf, 0x2f, 0x92, 0x4b, 0x7e, 0xbb, 0x94, 0x3c, 0x93, 0xe1, 0x32,
		0x11, 0x37, 0x4b, 0x1e, 0x2e, 0x22, 0x29, 0xbf, 0xa7, 0xd9, 0x65, 0x99,
		0xd4, 0xde, 0x79, 0xcc, 0xf2, 0x19, 0x0d, 0x5c, 0xfb, 0x3c, 0x60, 0xd4,
		0x76, 0x4f, 0xd8, 0x05, 0x35, 0x9e, 0xdc, 0x6f, 0xd0, 0xa9, 0xbb, 0x92,
		0x0c, 0x7a, 0x68, 0x34, 0x69, 0x33, 0x1a, 0x6b, 0xe4, 0x85, 0xea, 0xa5,
		0x75, 0x7b, 0xf1, 0x79, 0x24, 0xe2, 0x32, 0x33, 0xd9, 0xdd, 0xa8, 0xdc,
		0xdc, 0xea, 0x52, 0xe5, 0xcc, 0xeb, 0xd1, 0xd6, 0xeb, 0x14, 0x41, 0x12,
		0xcd, 0x79, 0x99, 0x1c, 0xef, 0xee, 0x52, 0xef, 0x6f, 0x35, 0x6a, 0xd2,
		0xe6, 0xc3, 0xde, 0x50, 0x23, 0x84, 0xe0, 0x1f, 0x77, 0x79, 0xf4, 0x29,
		0xae, 0xce, 0x2d, 0x1f, 0xa3, 0xd5, 0xe1, 0xbe, 0xf5, 0xd6, 0x61, 0xeb,
		0x52, 0xfd, 0x99, 0x21, 0x94, 0x39, 0x35, 0xb6, 0xeb, 0xb3, 0xf7, 0xcc,
		0xa3, 0x67, 0x9e, 0x7d, 0x6a, 0x79, 0x33, 0xfa, 0x81, 0xcd, 0xa8, 0x15,
		0xf8, 0x53, 0xdb, 0x55, 0xe5, 0xa7, 0xcc, 0xf5, 0x8f, 0x55, 0xc1, 0xfa,
		0x25, 0x28, 0xf5, 0xd9, 0x85, 0x4f, 0xdd, 0xa9, 0x5a, 0x81, 0xe3, 0x14,
		0xe2, 0x95, 0xc8, 0x64, 0x1e, 0x56, 0x72, 0x47, 0x8c, 0xa3, 0x46, 0xeb,
		0x8a, 0x2b, 0x07, 0x4b, 0x3a, 0xe2, 0xe2, 0x4b, 0x9a, 0xf0, 0xad, 0x62,
		0xf3, 0x23, 0xb7, 0x44, 0xdd, 0x3c, 0x1b, 0x4c, 0xc8, 0xd1, 0x11, 0x9e,
		0x95, 0x6e, 0xc8, 0x9b, 0x58, 0xfd, 0xfd, 0x43, 0xc9, 0x95, 0x95, 0xc9,
		0xe7, 0x76, 0x38, 0xdc, 0x70, 0xa8, 0x25, 0x1e, 0x16, 0x6f, 0x7d, 0xac,
		0x22, 0xf3, 0xfe, 0x35, 0xd6, 0xc8, 0xc1, 0x01, 0x7e, 0x60, 0xd5, 0x99,
		0x5c, 0x4a, 0x91, 0x26, 0xb2, 0xfe, 0x1c, 0x6c, 0xfa, 0x5c, 0xa7, 0x5b,
		0x56, 0xff, 0xb3, 0xcf, 0x61, 0x51, 0x52, 0xef, 0x7e, 0x7c, 0xe7, 0xaf,
		0xfc, 0xce, 0x78, 0xda, 0x0d, 0x7e, 0xbb, 0x10, 0x59, 0xe1, 0x55, 0xbb,
		0x4c, 0x37, 0x8b, 0xd9, 0xdc, 0x98, 0x6a, 0x00, 0x00, 0x7a, 0x03, 0xcc,
		0x3f, 0x00, 0xf4, 0x17, 0x98, 0x7f, 0x00, 0xe8, 0x2f, 0x7a, 0xf1, 0x80,
		0xf9, 0x07, 0x80, 0x5e, 0x02, 0xf3, 0x0f, 0x00, 0xfd, 0x05, 0xe6, 0x1f,
		0x00, 0xfa, 0xcb, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2e, 0xa4, 0x69,
		0x2b, 0x00, 0x1c, 0x00, 0x00,
	},
		"res/sqlite/deltaiota.db",
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
		}},
	}},
}}
