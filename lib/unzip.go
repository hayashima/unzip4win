package unzip4win

import (
	"errors"
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

func Unzip(zipPath string, config *Config, logger *zap.Logger) error {
	logger.Debug("Start unzip!")

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer func() {
		logger.Debug("Close")
		_ = reader.Close()
	}()

	startTime := time.Now()
	password, err := analyzePassword(reader.File[0], startTime, config, logger)
	if err != nil {
		return err
	}
	outputDir := outputDir(zipPath, config.Output)

	for _, f := range reader.File {
		err := save(f, password, outputDir, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func targetSpec(target time.Time, specs []SpecConfig, logger *zap.Logger) []SpecConfig {
	if specs == nil || len(specs) == 0 {
		logger.Debug("No match specs")
		return nil
	}
	if target.Sub(specs[0].StartDate) > 0 {
		logger.Debug("Match spec", zap.Any("spec", specs[0]))
		return specs
	}
	return targetSpec(target, specs[1:], logger)
}

func analyzePassword(f *zip.File, startDate time.Time, config *Config, logger *zap.Logger) (string, error) {
	if !f.IsEncrypted() {
		return "", nil
	}
	specs := config.Spec[:]
	for i := 0; i < config.Password.TryDays; i++ {
		targetDate := startDate.Add(time.Duration(-24*i) * time.Hour)
		specs := targetSpec(targetDate, specs, logger)
		if specs == nil {
			break
		}
		password := targetDate.Format(specs[0].Format)
		logger.Debug("try:", zap.String("password", password))
		f.SetPassword(password)
		if tryOpen(f) {
			logger.Debug("Match!", zap.String("password", password))
			return password, nil
		}
	}

	return "", errors.New("can't analyze password")
}

func tryOpen(f *zip.File) bool {
	r, err := f.Open()
	defer func() { _ = r.Close() }()
	if err != nil {
		return false
	}
	_, err = ioutil.ReadAll(r)
	return err == nil
}

func save(f *zip.File, password string, dest string, logger *zap.Logger) error {
	if f.IsEncrypted() {
		f.SetPassword(password)
	}
	r, err := f.Open()
	if err != nil {
		return nil
	}
	defer func() { _ = r.Close() }()

	path := filepath.Join(dest, f.Name)
	if f.FileInfo().IsDir() {
		logger.Debug("Create Dir", zap.String("dir", path))
		return os.MkdirAll(path, 0755)
	}

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	logger.Debug("Save",
		zap.String("file", path),
		zap.Int("size", len(buf)),
		zap.Any("mode", f.Mode()))
	err = ioutil.WriteFile(path, buf, f.Mode())
	return err
}

func outputDir(zipFile string, config OutputConfig) string {
	if config.SaveCurrent {
		return filepath.Dir(zipFile)
	}
	return config.OutputPath
}
