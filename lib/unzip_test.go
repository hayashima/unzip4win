package unzip4win

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestIsLookLikeZipFile(t *testing.T) {
	t.Run("is zip file", func(t *testing.T) {
		path := filepath.Join(".", "hoge", "fuga.zip")
		actual := IsLookLikeZipFile(path)
		if !actual {
			t.Error("expected: true, but actual is false")
		}
	})
	t.Run("is not zip file", func(t *testing.T) {
		path := filepath.Join(".", "hoge", "fuga.zip.txt")
		actual := IsLookLikeZipFile(path)
		if actual {
			t.Error("expected: false, but actual is true")
		}
	})
}

func TestTargetSpec(t *testing.T) {
	specs := []SpecConfig{
		{Format: "pass1", StartDate: createDate(2019, time.January, 1)},
		{Format: "pass2", StartDate: createDate(2018, time.January, 1)},
		{Format: "pass3", StartDate: createDate(2017, time.January, 1)},
	}
	t.Run("empty spec", func(t *testing.T) {
		target := createDate(2019, time.April, 1)
		actual := targetSpec(target, []SpecConfig{})
		if len(actual) != 0 {
			t.Errorf("expected size: 0, but actual size is %v\n actual is %v",
				len(actual), actual)
		}
	})

	t.Run("first is matched", func(t *testing.T) {
		target := createDate(2019, time.April, 1)
		actual := targetSpec(target, specs)
		if !reflect.DeepEqual(specs, actual) {
			t.Errorf("expected is %v, but actutal is %v", specs, actual)
		}
	})

	t.Run("first is NOT matched", func(t *testing.T) {
		target := createDate(2018, time.April, 1)
		actual := targetSpec(target, specs)
		if !reflect.DeepEqual(specs[1:], actual) {
			t.Errorf("expected is %v, but actutal is %v", specs[1:], actual)
		}
	})

	t.Run("nothing is matched", func(t *testing.T) {
		target := createDate(2016, time.December, 31)
		actual := targetSpec(target, specs)
		if len(actual) > 0 {
			t.Errorf("no match is expected, but actutal is %v", actual)
		}
	})
}

func TestOutputDir(t *testing.T) {
	t.Run("output current is true", func(t *testing.T) {
		zipPath := filepath.Join("..", "testdata", "zip", "test.zip")
		actual := outputDir(zipPath, OutputConfig{SaveCurrent: true})
		expected := filepath.Join("..", "testdata", "zip")
		if actual != expected {
			t.Errorf("expected is %v, but actual is %v", expected, actual)
		}
	})

	t.Run("output current is false", func(t *testing.T) {
		zipPath := filepath.Join("..", "testdata", "zip", "test.zip")
		expected := filepath.Join("..", "testdata", "zip", "output")
		actual := outputDir(zipPath, OutputConfig{SaveCurrent: false, OutputPath: expected})
		if actual != expected {
			t.Errorf("expected is %v, but actual is %v", expected, actual)
		}
	})
}

func createDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func TestUnzip(t *testing.T) {
	zipDir := filepath.Join("..", "testdata", "zip")
	outputDir := filepath.Join("..", "out", "zip")
	if exists(outputDir) {
		err := os.RemoveAll(outputDir)
		if err != nil {
			t.Fatal(err)
		}
	}
	config := Config{
		Output:   OutputConfig{SaveCurrent: false, OutputPath: outputDir},
		Password: PasswordConfig{TryDays: 10},
		Spec:     []SpecConfig{},
	}

	t.Run("utf-8 japanese file name", func(t *testing.T) {
		err := Unzip(filepath.Join(zipDir, "jp_file_utf8.zip"), &config)
		if err != nil {
			t.Fatal(err)
		}
		if !exists(filepath.Join(outputDir, "日本語ファイル名_utf8.txt")) {
			t.Error("wanted file `日本語ファイル名_utf8.txt` is not unzipped.")
		}
	})

	t.Run("utf-8 japanese dir and file names", func(t *testing.T) {
		err := Unzip(filepath.Join(zipDir, "jp_dir_utf8.zip"), &config)
		if err != nil {
			t.Fatal(err)
		}
		if !exists(filepath.Join(outputDir, "日本語ディレクトリ_utf8")) {
			t.Error("wanted dir `日本語ディレクトリ_utf8` is not unzipped.")
		}
		if !exists(filepath.Join(outputDir, "日本語ディレクトリ_utf8", "日本語ファイル名_utf8.txt")) {
			t.Error("wanted file `日本語ディレクトリ_utf8/日本語ファイル名_utf8.txt` is not unzipped.")
		}
	})

	t.Run("sjis japanese file name", func(t *testing.T) {
		err := Unzip(filepath.Join(zipDir, "jp_file_sjis.zip"), &config)
		if err != nil {
			t.Fatal(err)
		}
		if !exists(filepath.Join(outputDir, "日本語ファイル_sjis.txt")) {
			t.Error("wanted file `日本語ファイル_sjis.txt` is not unzipped.")
		}
	})

	t.Run("sjis japanese dir and file names", func(t *testing.T) {
		err := Unzip(filepath.Join(zipDir, "jp_dir_sjis.zip"), &config)
		if err != nil {
			t.Fatal(err)
		}
		if !exists(filepath.Join(outputDir, "日本語ディレクトリ_sjis")) {
			t.Error("wanted dir `日本語ディレクトリ_sjis` is not unzipped.")
		}
		if !exists(filepath.Join(outputDir, "日本語ディレクトリ_sjis", "日本語ファイル_sjis.txt")) {
			t.Error("wanted file `日本語ディレクトリ_sjis/日本語ファイル_sjis.txt` is not unzipped.")
		}
	})

	t.Run("password zip created by zip4win", func(t *testing.T) {
		config.Spec = []SpecConfig{{Format: "unzip", StartDate: createDate(2019, time.January, 1)}}
		err := Unzip(filepath.Join(zipDir, "jp_file_zip4win.zip"), &config)
		if err != nil {
			t.Fatal(err)
		}
		if !exists(filepath.Join(outputDir, "日本語ファイル名_zip4win.txt")) {
			t.Error("wanted file `日本語ファイル名_zip4win.txt` is not unzipped.")
		}
	})

	t.Run("directory is top in zip", func(t *testing.T) {
		config.Spec = []SpecConfig{{Format: "unzip", StartDate: createDate(2019, time.January, 1)}}
		err := Unzip(filepath.Join(zipDir, "dir_is_top.zip"), &config)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
