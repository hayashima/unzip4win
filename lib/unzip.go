package unzip4win

import (
	"errors"
	"github.com/yeka/zip"
	"go.uber.org/zap"
	"io/ioutil"
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
	password, nil := analyzePassword(reader.File[0], startTime, config, logger)

	for _, f := range reader.File {
		if f.IsEncrypted() {
			f.SetPassword(*password)
		}
		r, err := f.Open()
		if err != nil {
			return err
		}
		_, err = ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		logger.Debug("Open", zap.String("file", f.Name))
		_ = r.Close()
	}

	return nil
}

func targetSpec(target time.Time, specs []SpecConfig) []SpecConfig {
	if specs == nil || len(specs) == 0 {
		return nil
	}
	if target.Sub(specs[0].StartDate) > 0 {
		return specs
	}
	return targetSpec(target, specs[1:])
}

func analyzePassword(f *zip.File, startDate time.Time, config *Config, logger *zap.Logger) (*string, error) {
	if !f.IsEncrypted() {
		return nil, nil
	}
	specs := config.Spec[:]
	for i := 0; i < config.Password.TryDays; i++ {
		targetDate := startDate.Add(time.Duration(-24*i) * time.Hour)
		specs := targetSpec(targetDate, specs)
		if specs == nil {
			break
		}
		password := targetDate.Format(specs[0].Format)
		logger.Debug("try:", zap.String("password", password))
		f.SetPassword(password)
		if tryOpen(f) {
			logger.Debug("Match!", zap.String("password", password))
			return &password, nil
		}
	}

	return nil, errors.New("can't analyze password")
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
