package unzip4win

import (
	"github.com/saintfish/chardet"
	"go.uber.org/zap"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

var decoders = map[string]*encoding.Decoder{
	"Shift_JIS":   japanese.ShiftJIS.NewDecoder(),
	"EUC-JP":      japanese.EUCJP.NewDecoder(),
	"ISO-2022-JP": japanese.ISO2022JP.NewDecoder(),
}
var detector = chardet.NewTextDetector()

func decodeString(original string) (string, error) {
	detected, err := detector.DetectAll([]byte(original))
	if err != nil {
		return "", nil
	}
	debugLog("detected encodings", zap.Any("detected", detected))
	for _, e := range detected {
		decoder, contain := decoders[e.Charset]
		if contain {
			debugLog("matched encoding", zap.Any("encoding", e))
			return decoder.String(original)
		}
	}
	return original, nil
}
