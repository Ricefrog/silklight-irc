package main

import "fmt"

func formatString(line string, width int) []string {
	formatted := make([]string, 0, len(line)/width+1)

	for len(line) > width {
		formatted = append(formatted, line[:width])
		line = line[width:]
	}
	if len(line) > 0 {
		formatted = append(formatted, line)
	}

	return formatted
}

func printLines(lines []string) {
	for _, l := range lines {
		fmt.Println(l)
	}
}

func main() {
	//utils.RegTestStub()
	printLines(formatString("is this thing on?", 2))
	fmt.Println(len(formatString("is this thing on?", 2)))
	printLines(formatString("is this thing on?", 200))
}
