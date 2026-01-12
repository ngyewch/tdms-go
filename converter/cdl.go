package converter

import (
	"fmt"
	"os"
	"time"

	"github.com/ngyewch/tdms-go"
)

func ConvertToCDL(inputFile string, outputFile string) error {
	tdmsFile, err := tdms.OpenFile(inputFile)
	if err != nil {
		return err
	}
	defer func(tdmsFile *tdms.File) {
		_ = tdmsFile.Close()
	}(tdmsFile)

	datasetMap := make(map[string][]float64)
	channels := make([]*tdms.Node, 0)

	err = tdmsFile.ReadData(func(chunk tdms.Chunk) error {
		for _, channel := range chunk.Channels {
			values, exists := datasetMap[channel.Path]
			if !exists {
				channels = append(channels, channel.Node)
			}
			datasetMap[channel.Path] = append(values, channel.Samples...)
		}
		return nil
	})
	if err != nil {
		return err
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = f.WriteString("netcdf data {\n")
	if err != nil {
		return err
	}

	_, err = f.WriteString("\tdimensions:\n")
	if err != nil {
		return err
	}
	for _, channel := range channels {
		variableName := normalizeNetCDFIdentifier(channel.Name())
		values := datasetMap[channel.Path()]
		_, err = f.WriteString(fmt.Sprintf("\t\t%s = %d ;\n", variableName, len(values)))
		if err != nil {
			return err
		}
	}

	_, err = f.WriteString("\tvariables:\n")
	if err != nil {
		return err
	}

	for _, channel := range channels {
		variableName := normalizeNetCDFIdentifier(channel.Name())
		_, err = f.WriteString(fmt.Sprintf("\t\tdouble %s(%s) ;\n", variableName, variableName))
		if err != nil {
			return err
		}

		for propertyName, propertyValue := range channel.Properties().All() {
			convertedPropertyValue := func(propertyValue any) string {
				switch v := propertyValue.(type) {
				case time.Time:
					return v.Format(time.RFC3339)
				default:
					return fmt.Sprintf("%v", propertyValue)
				}
			}(propertyValue)
			_, err = f.WriteString(fmt.Sprintf("\t\t\t%s:%s = \"%s\" ;\n", variableName, propertyName, convertedPropertyValue))
			if err != nil {
				return err
			}
		}
	}

	_, err = f.WriteString("}\n")
	if err != nil {
		return err
	}

	return nil
}

func normalizeNetCDFIdentifier(s string) string {
	var s2 string
	for _, c := range s {
		if isValidNetCDFIdentifierChar(c) {
			s2 += string(c)
		} else {
			s2 += "_"
		}
	}
	return s2
}

func isValidNetCDFIdentifierChar(c rune) bool {
	return ((c >= 'a') && (c <= 'z')) ||
		((c >= 'A') && (c <= 'Z')) ||
		((c >= '0') && (c <= '9')) ||
		(c == '!') ||
		(c == '#') ||
		(c == '$') ||
		(c == '%') ||
		(c == '&') ||
		(c == '*') ||
		(c == ':') ||
		(c == ';') ||
		(c == '<') ||
		(c == '=') ||
		(c == '>') ||
		(c == '?') ||
		(c == '/') ||
		(c == '^') ||
		(c == '|') ||
		(c == '~') ||
		(c == '_') ||
		(c == '.') ||
		(c == '@') ||
		(c == '+') ||
		(c == '-')
}

/*
func isValidNetCDFIdentifierChar(c rune) bool {
	return ((c >= 'a') && (c <= 'z')) ||
		((c >= 'A') && (c <= 'Z')) ||
		((c >= '0') && (c <= '9')) ||
		(c == '!') ||
		(c == '#') ||
		(c == '$') ||
		(c == '%') ||
		(c == '&') ||
		(c == '*') ||
		(c == ':') ||
		(c == ';') ||
		(c == '<') ||
		(c == '=') ||
		(c == '>') ||
		(c == '?') ||
		(c == '/') ||
		(c == '^') ||
		(c == '|') ||
		(c == '~') ||
		(c == '_') ||
		(c == '.') ||
		(c == '@') ||
		(c == '+') ||
		(c == '-')
}
*/
