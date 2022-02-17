package fs

import (
	"encoding/json"

	fsitem "github.com/DevbookHQ/template-init/pkg/fs/item"
)

type DirItem struct {
	path string
}

func (di *DirItem) Path() string      { return di.path }
func (di *DirItem) Type() fsitem.Type { return fsitem.DIR }

func (di *DirItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ItemType fsitem.Type `json:"type"`
		Path     string      `json:"path"`
	}{
		ItemType: di.Type(),
		Path:     di.Path(),
	})
}

func (di *DirItem) UnmarshalJSON(b []byte) error {
	var item struct {
		Path string `json:"path"`
	}
	err := json.Unmarshal(b, &item)
	if err != nil {
		return err
	}

	di.path = item.Path
	return nil
}
