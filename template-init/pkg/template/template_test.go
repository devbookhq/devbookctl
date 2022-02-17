package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"sort"
	"testing"

	"github.com/DevbookHQ/template-init/pkg/environment"
)

var (
	testDir        = "./test/"
	pkgJSONPath    = testDir + "package.json"
	pkgJSONContent = `
{
  "name": "runner-template-nextjs-v11",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "next": "11.1.2",
    "react": "17.0.2",
    "react-dom": "17.0.2",
		"@devbookhq/splitter": "^1.0.0"
  },
  "devDependencies": {
    "@types/react": "17.0.21",
    "eslint": "7.32.0",
    "eslint-config-next": "11.1.2",
    "typescript": "4.4.3"
  }
}
	`
)

func TestInstallingPackages(t *testing.T) {
	// Cleanup possible left overs from the previous session.
	cleanup()
	// Cleanup after self.
	defer cleanup()

	mkTestDir()
	mkPackageJSON()

	testTempl, err := FromEnvironment(&environment.Environment{
		RUNNER_SOCKET_PATH: testDir + "test_runner.socket",
		ROOT_DIR:           testDir,
		CODE_CELLS_DIR:     testDir + "cc/",
		START_CMD:          "echo hello",
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		toInstall         []string
		expInstalledStart []string
		expDiff           []string
		expInstalledEnd   []string
	}{
		{
			toInstall:         []string{"next", "react", "express", "jquery"},
			expInstalledStart: []string{"next", "react", "react-dom", "@devbookhq/splitter"},
			expDiff:           []string{"jquery", "express"},
			expInstalledEnd:   []string{"next", "react", "react-dom", "@devbookhq/splitter", "jquery", "express"},
		},
		{
			toInstall:         []string{},
			expInstalledStart: []string{"next", "react", "react-dom", "@devbookhq/splitter"},
			expDiff:           []string{},
			expInstalledEnd:   []string{"next", "react", "react-dom", "@devbookhq/splitter"},
		},
		{
			toInstall:         []string{"next", "react", "react-dom", "@devbookhq/splitter"},
			expInstalledStart: []string{"next", "react", "react-dom", "@devbookhq/splitter"},
			expDiff:           []string{},
			expInstalledEnd:   []string{"next", "react", "react-dom", "@devbookhq/splitter"},
		},
	}

	for _, test := range tests {
		t.Logf("++++++++++++++++++++++")
		t.Logf("++++++++++++++++++++++")
		cleanup()
		mkTestDir()
		mkPackageJSON()

		// Test installed packages in the beginning.
		err, installedStart := testTempl.installedPackages()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("installedStart: %s", installedStart)
		t.Logf("expInstalledStart: %s", test.expInstalledStart)
		compareArrs(t, test.expInstalledStart, installedStart)
		t.Logf("======== Ok")

		// Test diffing the installed packages and packages we want to install.
		diff := testTempl.diff(installedStart, test.toInstall)
		t.Logf("installedStart: %s", installedStart)
		t.Logf("toInstall: %s", test.toInstall)
		t.Logf("diff: %s", diff)
		t.Logf("expDiff: %s", test.expDiff)
		compareArrs(t, test.expDiff, diff)
		t.Logf("======== Ok")

		// Test newly installed packages.
		out, err := testTempl.InstallPackages(test.toInstall)
		if err != nil {
			t.Fatalf("%s: %s", string(out), err)
		}
		err, installedEnd := testTempl.installedPackages()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("installedEnd: %s", installedEnd)
		t.Logf("expInstalledEnd: %s", test.expInstalledEnd)
		compareArrs(t, test.expInstalledEnd, installedEnd)
		t.Logf("======== Ok")
	}
}

func TestParsingPackageJSON(t *testing.T) {
	// Cleanup possible left overs from the previous session.
	cleanup()
	// Cleanup after self.
	defer cleanup()

	mkTestDir()
	mkPackageJSON()

	testTempl, err := FromEnvironment(&environment.Environment{
		RUNNER_SOCKET_PATH: testDir + "test_runner.socket",
		ROOT_DIR:           testDir,
		CODE_CELLS_DIR:     testDir + "cc/",
		START_CMD:          "echo hello",
	})
	if err != nil {
		t.Fatal(err)
	}

	err, installedPkgs := testTempl.installedPackages()
	if err != nil {
		t.Fatal(err)
	}
	expectedPkgs := []string{"next", "react", "react-dom", "@devbookhq/splitter"}

	t.Logf("expectedPkgs: %s", expectedPkgs)
	t.Logf("installedPkgs: %s", installedPkgs)

	compareArrs(t, expectedPkgs, installedPkgs)
}

func mkTestDir() {
	if err := os.MkdirAll(testDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to make test dir: %s", err))
	}
}

func mkPackageJSON() {
	if err := ioutil.WriteFile(pkgJSONPath, []byte(pkgJSONContent), 0755); err != nil {
		panic(fmt.Sprintf("failed to create package.json: %s", err))
	}
}

func compareArrs(t *testing.T, expected, got []string) {
	if len(expected) != len(got) {
		t.Log(string(debug.Stack()))
		t.Fatal("len(expected) and len(got) don't match")
	}

	sort.Strings(expected)
	sort.Strings(got)

	for idx, e := range expected {
		g := got[idx]
		if e != g {
			t.Log(string(debug.Stack()))
			t.Fatalf("expected and got don't match. expected=%s got=%s", e, g)
		}
	}
}

func cleanup() {
	err := os.RemoveAll(testDir)
	if os.IsNotExist(err) {
		// Ignore if the file doesn't exist.
		return
	}
	if err != nil {
		panic(err)
	}
}
