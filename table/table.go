package table

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	lengthToPrint     = 150
	halfLengthToPrint = lengthToPrint / 2
	lengthOfTime      = 15

	// ANSI Scape Color
	black      = "\u001b[1;30m"
	red        = "\u001b[1;31m"
	green      = "\u001b[1;32m"
	yellow     = "\u001b[1;33m"
	blue       = "\u001b[1;34m"
	magenta    = "\u001b[1;35m"
	cyan       = "\u001b[1;36m"
	white      = "\u001b[1;37m"
	colorReset = "\033[0m"
)

func Distribution(name string, minValue float64, q1 float64, medianValue float64, meanValue float64, q3 float64, maxValue float64, total int) {
	var message = fmt.Sprint(name+" Distribution",
		" | Min: ", minValue,
		" | Q1: ", q1,
		" | Median: ", medianValue,
		" | Mean: ", meanValue,
		" | Q3: ", q3,
		" | Max: ", maxValue,
		" | Total: ", total)

	var size = lengthToPrint - utf8.RuneCountInString(message) - 1
	var symbol = strings.Repeat(" ", size)

	fmt.Print(string(blue), "| ", string(colorReset),
		string(green), name+" Distribution", string(colorReset),
		string(blue), " | ", string(colorReset),
		string(green), fmt.Sprint("Min: ", minValue), string(colorReset),
		string(blue), " | ", string(colorReset),
		string(green), fmt.Sprint("Q1: ", q1), string(colorReset),
		string(blue), " | ", string(colorReset),
		string(green), fmt.Sprint("Median: ", medianValue), string(colorReset),
		string(blue), " | ", string(colorReset),
		string(green), fmt.Sprint("Mean: ", meanValue), string(colorReset),
		string(blue), " | ", string(colorReset),
		string(green), fmt.Sprint("Q3: ", q3), string(colorReset),
		string(blue), " | ", string(colorReset),
		string(green), fmt.Sprint("Max: ", maxValue), string(colorReset),
		string(blue), " | ", string(colorReset),
		string(green), fmt.Sprint("Total: ", total), string(colorReset),
		symbol, string(blue), "|\n", string(colorReset))

	NewDivisor()
}

func LoopStepsPrint(steps ...float64) {
	var message = "Step Number: "
	for i := range steps {
		if i < 1 {
			message += fmt.Sprint(steps[i])
			continue
		}
		message = message + " - " + fmt.Sprint(steps[i])
	}
	if utf8.RuneCountInString(message)%2 != 0 {
		message = " " + message
	}
	var size = (lengthToPrint - utf8.RuneCountInString(message)) / 2
	var symbol = strings.Repeat(" ", size) + message + strings.Repeat(" ", size)
	fmt.Print(string(blue), "|", string(colorReset), string(red), symbol, string(colorReset), string(blue), "|\n", string(colorReset))
	NewDivisor()
}

func StartHead() {
	NewDivisor()
	var size = (lengthToPrint - utf8.RuneCountInString("Starting")) / 2
	var symbol = strings.Repeat(" ", size) + "Starting" + strings.Repeat(" ", size)
	fmt.Print(string(blue), "|", string(colorReset), string(red), symbol, string(colorReset), string(blue), "|\n", string(colorReset))
	NewDivisor()
}

func AnalysisHead() {
	var size = (lengthToPrint - utf8.RuneCountInString("Analysis")) / 2
	var symbol = strings.Repeat(" ", size) + "Analysis" + strings.Repeat(" ", size)
	fmt.Print(string(blue), "|", string(colorReset), string(red), symbol, string(colorReset), string(blue), "|\n", string(colorReset))
	NewDivisor()
}

func AnalysisBody(message string, color string) {
	var size = lengthToPrint - utf8.RuneCountInString(message) - 1
	var symbol = message + strings.Repeat(" ", size)
	if color == "white" {
		fmt.Print(string(blue), "| ", string(colorReset), string(white), symbol, string(colorReset), string(blue), "|\n", string(colorReset))
	} else if color == "green" {
		fmt.Print(string(blue), "| ", string(colorReset), string(green), symbol, string(colorReset), string(blue), "|\n", string(colorReset))
	} else if color == "red" {
		fmt.Print(string(blue), "| ", string(colorReset), string(red), symbol, string(colorReset), string(blue), "|\n", string(colorReset))
	}
	NewDivisor()
}

func StartBody(message string, elapsed time.Duration) {
	var strTime = fmt.Sprint(elapsed)
	var finalSize = lengthOfTime - utf8.RuneCountInString(strTime)
	var finalSpaces = strings.Repeat(" ", finalSize)
	var size = lengthToPrint - (utf8.RuneCountInString("| Time: |") + utf8.RuneCountInString(message) + utf8.RuneCountInString(strTime) + finalSize)
	var spaces = strings.Repeat(" ", size)
	fmt.Print(string(blue), "| ", string(colorReset),
		string(white), message, spaces, string(colorReset),
		string(blue), "| ", string(colorReset),
		string(yellow), "Time: ", elapsed, finalSpaces, string(yellow),
		string(blue), "|\n", string(colorReset))
	NewDivisor()
}

func NewDivisor() {
	var symbol = strings.Repeat("-", halfLengthToPrint)
	var toPrint = "+" + symbol + symbol + "+\n"
	fmt.Print(blue, toPrint, colorReset)
}
