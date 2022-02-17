package fs

import (
	"encoding/json"

	fsitem "github.com/DevbookHQ/template-init/pkg/fs/item"
)

type FileItem struct {
	path string
}

func (fi *FileItem) Path() string      { return fi.path }
func (fi *FileItem) Type() fsitem.Type { return fsitem.FILE }

func (fi *FileItem) MarshalJSON() ([]byte, error) {
	//buffer := bytes.NewBufferString("{")

	//path := fmt.Sprintf("\"path\":\"%s\"", fi.Path())
	//buffer.WriteString(path)

	//itemtype := fmt.Sprintf("\"type\":\"%s\"", fi.Type())
	//buffer.WriteString(item)

	//buffer.WriteString("}")
	//return buffer, nil

	return json.Marshal(struct {
		ItemType fsitem.Type `json:"type"`
		Path     string      `json:"path"`
	}{
		ItemType: fi.Type(),
		Path:     fi.Path(),
	})
}

func (fi *FileItem) UnmarshalJSON(b []byte) error {
	var item struct {
		Path string `json:"path"`
	}
	err := json.Unmarshal(b, &item)
	if err != nil {
		return err
	}

	fi.path = item.Path
	return nil
}
