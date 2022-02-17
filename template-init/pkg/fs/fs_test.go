package fs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

//func TestRemoveFile(t *testing.T) {
//	testDir := tempMkdir(t)
//	defer os.RemoveAll(testDir)
//	t.Logf("testDir: %s\n", testDir)
//
//	testEvents := []RemoveEvent{
//		{Path: "/test"},
//	}
//
//	fs, err := New(testDir)
//	if err != nil {
//		t.Fatalf("failed to create new Filesystem: %s", err)
//	}
//
//	go func() {
//		for _, te := range testEvents {
//			t.Logf("remove event: %+v\n", te)
//			fs.Do(&te)
//			p := filepath.Join(fs.RootDir, te.Path)
//			checkFileAndContent(t, p, te.Content)
//		}
//		fs.Errors <- nil
//	}()
//
//}

func TestMarshalFileItem(t *testing.T) {
	fi := &FileItem{"/test/path"}

	b, err := json.Marshal(fi)
	if err != nil {
		t.Fatalf("Failed to marshal FileItem: %s", err)
	}
	t.Logf("Marshaled FileItem: %s", string(b))

	expectedB, err := json.Marshal(struct {
		Type string `json:"type"`
		Path string `json:"path"`
	}{"File", "/test/path"})

	res := bytes.Compare(b, expectedB)
	if res != 0 {
		t.Fatal("Marshaled bytes don't match expected bytes")
	}
}

func TestUnmarshalFileItem(t *testing.T) {
	j := `{"type":"File","path":"/test/path"}`

	var fi FileItem
	if err := json.Unmarshal([]byte(j), &fi); err != nil {
		t.Fatalf("Failed to unmarshal FileItem: %s", err)
	}
	t.Logf("Unmarshaled FileItem: %+v\n", fi)

	expected := &FileItem{"/test/path"}
	if expected.Type() != fi.Type() {
		t.Fatalf("Type() calls don't match. expected=%s, got=%s", expected.Type(), fi.Type())
	}
	if expected.Path() != fi.Path() {
		t.Fatalf("Path() calls don't match. expected=%s, got=%s", expected.Path(), fi.Path())
	}
}

func TestMarshalDirItem(t *testing.T) {
	fi := &DirItem{"/test/path"}

	b, err := json.Marshal(fi)
	if err != nil {
		t.Fatalf("Failed to marshal FileItem: %s", err)
	}
	t.Logf("Marshaled DirItem: %s", string(b))

	expectedB, err := json.Marshal(struct {
		Type string `json:"type"`
		Path string `json:"path"`
	}{"Directory", "/test/path"})

	res := bytes.Compare(b, expectedB)
	if res != 0 {
		t.Fatal("Marshaled bytes don't match expected bytes")
	}
}

func TestUnmarshalDirItem(t *testing.T) {
	j := `{"type":"Dir","path":"/test/path"}`

	var di DirItem
	if err := json.Unmarshal([]byte(j), &di); err != nil {
		t.Fatalf("Failed to unmarshal FileItem: %s", err)
	}
	t.Logf("Unmarshaled DirItem: %+v\n", di)

	expected := &DirItem{"/test/path"}
	if expected.Type() != di.Type() {
		t.Fatalf("Type() calls don't match. expected=%s, got=%s", expected.Type(), di.Type())
	}
	if expected.Path() != di.Path() {
		t.Fatalf("Path() calls don't match. expected=%s, got=%s", expected.Path(), di.Path())
	}
}

//func TestWatchingEvents(t *testing.T) {
//	testDir := tempMkdir(t)
//	defer os.RemoveAll(testDir)
//	t.Logf("testDir: %s\n", testDir)
//
//	fs, err := New(testDir)
//	if err != nil {
//		t.Fatalf("Failed to create new Filesystem: %s", err)
//	}
//
//	content, err := fs.ListDir("/Users/vasekmlejnsky/Developer/node_project/")
//	if err != nil {
//		t.Fatalf("Failed to list dir: %s", err)
//	}
//
//	for _, i := range content {
//		t.Logf("Type: %s, Path: %s", i.Type(), i.Path())
//	}
//
//	done := make(chan struct{})
//	<-done
//}

//func TestCreateFile(t *testing.T) {
//	testDir := tempMkdir(t)
//	defer os.RemoveAll(testDir)
//	t.Logf("testDir: %s\n", testDir)
//
//	testEvents := []CreateEvent{
//		//{Path: "/test", Content: "hello"},
//		{Path: "/test"},
//	}
//
//	fs, err := New(testDir)
//	if err != nil {
//		t.Fatalf("failed to create new Filesystem: %s", err)
//	}
//
//	go func() {
//		for _, te := range testEvents {
//			t.Logf("create event: %+v\n", te)
//			fs.Do(&te)
//			p := filepath.Join(fs.RootDir, te.Path)
//			//checkFileAndContent(t, p, te.Content)
//			checkFileAndContent(t, p, "")
//		}
//		fs.Errors <- nil
//	}()
//
//	err = <-fs.Errors
//	if err != nil {
//		t.Fatalf("filesytem error: %s", err)
//	}
//}
//
//func checkFileAndContent(t *testing.T, path, content string) {
//	dat, err := os.ReadFile(path)
//	if err != nil {
//		t.Fatalf("failed to read file: path=%s, err=%s", path, err)
//	}
//
//	if string(dat) != content {
//		t.Fatalf(
//			"file content doesn't match expected value: path=%s, got=%s, expected=%s",
//			path,
//			string(dat),
//			content,
//		)
//	}
//}
//
func tempMkdir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "tinit_fs")
	if err != nil {
		t.Fatalf("failed to create test directory: %s", err)
	}
	return dir
}
