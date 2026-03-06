package test

import (
	"escobra/cmd"
	"os"
	"path/filepath"
	"testing"
)

func TestGetNameJsonCase(t *testing.T) {
	tmpDir := t.TempDir()
	file2 := filepath.Join(tmpDir, "test-index-pupis.json")
	err := os.WriteFile(file2, []byte(`{"test": true}`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Stat(file2)
	if err != nil {
		t.Fatal(err)
	}

	exp := "test-index-pupis"
	output := cmd.GetName(f)

	if exp != output {
		t.Errorf("expected %s, got %s", exp, output)
	}

}
func TestGetNameDirCase(t *testing.T) {
	tmpDir := t.TempDir()
	dir := filepath.Join(tmpDir, "test-index-log")
	err := os.Mkdir(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}

	exp := "test-index-log"
	output := cmd.GetName(f)

	if exp != output {
		t.Errorf("expected %s, got %s", exp, output)
	}

}
