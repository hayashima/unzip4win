package unzip4win

import (
	"github.com/pkg/errors"
	"github.com/yeka/zip"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func IsLookLikeZipFile(zipPath string) bool {
	return filepath.Ext(zipPath) == ".zip"
}

func Unzip(zipPath string, config *Config) error {
	debugLog("Start unzip!")

	outputDir := outputDir(zipPath, config.Output)
	err := createDir(outputDir)
	if err != nil {
		return errors.WithStack(err)
	}

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		debugLog("Close")
		_ = reader.Close()
	}()

	var password string
	for _, f := range reader.File {
		password, err = unzipFile(f, password, outputDir, config)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func createDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func targetSpec(target time.Time, specs []SpecConfig) []SpecConfig {
	if specs == nil || len(specs) == 0 {
		debugLog("No match specs")
		return nil
	}
	if target.Sub(specs[0].StartDate) > 0 {
		debugLog("Match spec", zap.Any("spec", specs[0]))
		return specs
	}
	return targetSpec(target, specs[1:])
}

func analyzePassword(f *zip.File, startDate time.Time, config *Config) (string, error) {
	if !f.IsEncrypted() {
		return "", nil
	}
	specs := config.Spec[:]
	for i := 0; i < config.Password.TryDays; i++ {
		targetDate := startDate.Add(time.Duration(-24*i) * time.Hour)
		specs := targetSpec(targetDate, specs)
		if specs == nil {
			break
		}
		password := targetDate.Format(specs[0].Format)
		debugLog("try:", zap.String("password", password))
		if tryOpen(f, password) {
			debugLog("Match!", zap.String("password", password))
			return password, nil
		}
	}

	return "", errors.New("can't analyze password")
}

func tryOpen(f *zip.File, password string) bool {
	f.SetPassword(password)
	r, err := f.Open()
	defer func() { _ = r.Close() }()
	if err != nil {
		return false
	}
	_, err = ioutil.ReadAll(r)
	return err == nil
}

func unzipFile(f *zip.File, password string, outputDir string, config *Config) (newPassword string, err error) {
	decodedName, err := decodeString(f.Name)
	if err != nil {
		return "", errors.WithStack(err)
	}

	path := filepath.Join(outputDir, decodedName)
	if f.FileInfo().IsDir() {
		debugLog("Create Dir", zap.String("dir", path))
		return password, createDir(path)
	}

	if f.IsEncrypted() {
		if tryOpen(f, password) {
			newPassword = password
		} else {
			newPassword, err = analyzePassword(f, time.Now(), config)
			if err != nil {
				return "", errors.WithStack(err)
			}
		}
	}

	r, err := f.Open()
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer func() { _ = r.Close() }()

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.WithStack(err)
	}
	debugLog("Save",
		zap.String("file", path),
		zap.Int("size", len(buf)),
		zap.Any("mode", f.Mode()))
	err = ioutil.WriteFile(path, buf, f.Mode())

	return newPassword, errors.WithStack(err)
}

func outputDir(zipFile string, config OutputConfig) string {
	if config.SaveCurrent {
		return filepath.Dir(zipFile)
	}
	return config.OutputPath
}
