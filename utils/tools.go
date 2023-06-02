package utils

import (
	"github.com/qmaru/minitools"
)

// DataSuite
var DataSuite *minitools.DataSuiteBasic

// FileSuite
var FileSuite *minitools.FileSuiteBasic

// TimeSuite
var TimeSuite *minitools.TimeSuiteBasic

func init() {
	DataSuite = minitools.DataSuite()
	FileSuite = minitools.FileSuite()
	TimeSuite = minitools.TimeSuite()
}
