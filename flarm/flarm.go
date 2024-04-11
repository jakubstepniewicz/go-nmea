// Package flarm parses FLARM NMEA sentences.
//
// See https://www.flarm.com/wp-content/uploads/man/FTD-012-Data-Port-Interface-Control-Document-ICD.pdf.
package flarm

import "github.com/twpayne/go-nmea"

var sentenceParserMap = nmea.SentenceParserMap{
	"PFLAA": nmea.MakeSentenceParser(ParsePFLAA),
	"PFLAE": ParsePFLAE,
	"PFLAU": nmea.MakeSentenceParser(ParsePFLAU),
	"PFLAV": ParsePFLAV,
}

func SentenceParserFunc(addr string) nmea.SentenceParser {
	return sentenceParserMap[addr]
}