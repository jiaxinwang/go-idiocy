package common

import (
	"fmt"
	"strings"

	"github.com/bregydoc/gtranslate"
	"github.com/sirupsen/logrus"
)

// TranslateMultiLines ...
func TranslateMultiLines(lines []string, from, to string) []string {
	translated := make([]string, len(lines))
	count, last := 0, 0
	sep := fmt.Sprintf("\n")
	for k, v := range lines {
		count += len(v)
		if count+len(v) > 4500 || k == len(lines)-1 {
			merged := strings.Join(lines[last:k+1], sep)
			ret, err := Translate(merged, from, to)
			if err != nil {
				logrus.WithError(err).Error()
			}
			translatedLines := strings.Split(ret, sep)
			for i := 0; i < k-last+1; i++ {
				translated[last+i] = translatedLines[i]
			}
			last = k
			count = 0
		}
	}
	return translated

}

// Translate string
func Translate(line string, from, to string) (translated string, err error) {
	translated, err = gtranslate.TranslateWithParams(
		line,
		gtranslate.TranslationParams{
			From: from,
			To:   to,
		},
	)
	if err != nil {
		return line, err
	}
	return translated, nil

}
