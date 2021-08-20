package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func setupTf(t *testing.T, dir string) {
	for _, file := range []string{"./terraform.d", "./providers.tf"} {
		abs, err := filepath.Abs(file)
		if err != nil {
			t.Fatal(err)
		}
		os.Symlink(abs, fmt.Sprintf("%s/%s", dir, file))
	}
}
func createFs(t *testing.T, files map[string]string) (string, func()) {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		t.Fatal(err)
	}

	for name, content := range files {
		err = ioutil.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	setupTf(t, dir)

	return dir, func() {
		os.RemoveAll(dir)
	}
}
