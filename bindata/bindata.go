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
		0x3f, 0x4b, 0xc3, 0x40, 0x18, 0xc6, 0xef, 0x92, 0x62, 0x31, 0xe0, 0xe4,
		0x90, 0xf5, 0xe5, 0xa6, 0x86, 0xd6, 0xc9, 0xc1, 0xa9, 0x60, 0xd4, 0x43,
		0x0e, 0xdb, 0x6b, 0x1b, 0x2f, 0xd0, 0x4c, 0x21, 0xda, 0x13, 0x22, 0x6d,
		0x4a, 0x9b, 0x04, 0x5c, 0xfd, 0x04, 0x7e, 0x36, 0xbf, 0x84, 0x9f, 0xc3,
		0x4b, 0x2c, 0x2d, 0x0d, 0x82, 0x93, 0x53, 0xde, 0x1f, 0xe4, 0x92, 0xe7,
		0xde, 0xbf, 0x19, 0x9e, 0xc7, 0xd9, 0x28, 0x2d, 0x34, 0xbc, 0xac, 0xb7,
		0xab, 0xa4, 0x80, 0x4b, 0xd2, 0x21, 0x94, 0x92, 0x6b, 0x00, 0x42, 0x08,
		0x35, 0x4f, 0x87, 0x1c, 0xb0, 0x1a, 0x9a, 0x92, 0xbf, 0xa1, 0xe4, 0xe2,
		0xcb, 0x3a, 0x33, 0x1f, 0xb6, 0xf5, 0x49, 0xec, 0x0f, 0x7b, 0x68, 0x5e,
		0x08, 0xf2, 0x7f, 0xbc, 0xda, 0x5d, 0xf7, 0xca, 0xa5, 0xef, 0xe7, 0x69,
		0xb6, 0xd0, 0x6f, 0x65, 0xae, 0xb7, 0x79, 0x5c, 0x66, 0xe9, 0xa6, 0xd4,
		0x71, 0x25, 0xb2, 0x64, 0xa5, 0xeb, 0xcb, 0xce, 0x6d, 0xc0, 0x7d, 0xc5,
		0x21, 0x94, 0x62, 0x16, 0x72, 0x10, 0xf2, 0x8e, 0xcf, 0x81, 0xfd, 0x9a,
		0xcf, 0x60, 0x22, 0x77, 0x21, 0x06, 0x3d, 0xb6, 0xbf, 0xf6, 0xa6, 0xd6,
		0x89, 0xdb, 0xef, 0xd3, 0xa8, 0x48, 0x9e, 0x96, 0x3a, 0xdf, 0x2c, 0x8d,
		0x93, 0xe2, 0x5c, 0x9b, 0xd2, 0xec, 0xb9, 0x29, 0xed, 0xdd, 0x38, 0xe5,
		0xdf, 0x8c, 0x38, 0x34, 0x82, 0xbd, 0xaa, 0xdd, 0xc0, 0x28, 0x6f, 0x4d,
		0xbb, 0xae, 0x6b, 0xb6, 0x1f, 0xd6, 0x2d, 0xeb, 0x91, 0xf5, 0x61, 0x1d,
		0x95, 0xef, 0x77, 0x71, 0x4e, 0x59, 0xba, 0x60, 0xf0, 0x83, 0x90, 0x8a,
		0xdf, 0xf3, 0x00, 0xa6, 0x81, 0x18, 0xfb, 0x41, 0x04, 0x0f, 0x3c, 0x02,
		0x3f, 0x54, 0x13, 0x21, 0x4d, 0xf1, 0x98, 0x4b, 0x35, 0x30, 0xe9, 0x87,
		0x7f, 0x52, 0x7c, 0xae, 0x1c, 0xaf, 0xf2, 0xe6, 0x91, 0xab, 0x11, 0x04,
		0x69, 0x0d, 0xe8, 0x7f, 0x04, 0x69, 0x2f, 0x4e, 0x75, 0xa0, 0xff, 0x11,
		0xa4, 0x95, 0x7c, 0x07, 0x00, 0x00, 0xff, 0xff, 0xd2, 0x64, 0x9a, 0xdb,
		0x00, 0x10, 0x00, 0x00,
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
