package test

import (
	"escobra/cmd"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckDirOrFile(t *testing.T) {

	p := "/ec/local/home/nikanni/my-programming/app-workspace/es-cli/cmd"

	fileInfo, _ := os.Stat(p)

	fmt.Println(fileInfo.IsDir())

}

func TestFileName(t *testing.T) {
	fileName := "/ec/local/home/nikanni/my-programming/app-workspace/es-cli/cmd"

	// fmt.Println(filepath.)

	fileInfo, _ := os.Stat(fileName)
	fmt.Printf("%+v\n", fileInfo)

	fmt.Println(filepath.Base(fileInfo.Name()))

}

//

// func TestGetName_File(t *testing.T) {
// 	tmpFile, err := os.CreateTemp("", "mappings.json")
// 	if err != nil {
// 		t.Fatalf("failed to create temp file: %v", err)
// 	}
// 	defer os.Remove(tmpFile.Name())

// 	expected := "example"
// 	result := cmd.GetName(info)

// 	if result != expected {
// 		t.Errorf("expected %s, got %s", expected, result)
// 	}
// }

func TestGetNameJsonCase(t *testing.T) {
	tmpDir := t.TempDir()
	file2 := filepath.Join(tmpDir, "test-index-pupis.json")
	err := os.WriteFile(file2, []byte(`{"test": true}`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	f2, err := os.Stat(file2)
	if err != nil {
		t.Fatal(err)
	}

	exp2 := "test-index-pupis"
	output2 := cmd.GetName(f2)

	if exp2 != output2 {
		t.Errorf("expected %s, got %s", exp2, output2)
	}

}
func TestGetNameDirCase(t *testing.T) {
	tmpDir := t.TempDir()
	dir1 := filepath.Join(tmpDir, "test-index-log")
	err := os.Mkdir(dir1, 0755)
	if err != nil {
		t.Fatal(err)
	}

	f1, err := os.Stat(dir1)
	if err != nil {
		t.Fatal(err)
	}

	exp1 := "test-index-log"
	output1 := cmd.GetName(f1)

	if exp1 != output1 {
		t.Errorf("expected %s, got %s", exp1, output1)
	}

}

// func TestGetName_Directory(t *testing.T) {
// 	tmpDir, err := os.MkdirTemp("", "exampledir")
// 	if err != nil {
// 		t.Fatalf("failed to create temp dir: %v", err)
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	info, err := os.Stat(tmpDir)
// 	if err != nil {
// 		t.Fatalf("failed to stat dir: %v", err)
// 	}

// 	expected := filepath.Base(tmpDir)
// 	result := cmd.GetName(info)

// 	if result != expected {
// 		t.Errorf("expected %s, got %s", expected, result)
// 	}
// }

// tmpDir, err := os.MkdirTemp("", "test-dir")
// if err != nil {
// 	t.Fatal(err)
// }

// tmpFileMappings, err := os.CreateTemp(tmpDir, "mappings.json*")
// // tmpFileSettings, err := os.CreateTemp(tmpDir, "settings.json")

// // fmt.Printf("%+v\n", tmpDir)
// fmt.Printf("%+v\n", tmpFileMappings.Name())
// // fmt.Printf("%+v", tmpFileSettings)

// tmpDir := t.TempDir()
// dir := filepath.Join(tmpDir, "test-index")
// fmt.Println(dir)

// filePath := filepath.Join(dir, "test-index-ms.json")

// fileInfo, err := os.Stat(filePath)
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf("%+v", fileInfo)
// // output := cmd.GetName(filePath)
