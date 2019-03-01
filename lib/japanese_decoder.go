package unzip4win

import (
	"github.com/saintfish/chardet"
	"go.uber.org/zap"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/unicode/norm"
)

var decoders = map[string]*encoding.Decoder{
	"Shift_JIS": japanese.ShiftJIS.NewDecoder(),
	// now, support file name which is encoded with SJIS or UTF-8 only.
	//"EUC-JP":      japanese.EUCJP.NewDecoder(),
	//"ISO-2022-JP": japanese.ISO2022JP.NewDecoder(),
	"UTF-8": nil,
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
			if decoder == nil {
				break
			} else {
				return decoder.String(original)
			}
		}
	}
	return string(norm.NFC.Bytes([]byte(original))), nil
}
