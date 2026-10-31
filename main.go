package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	AFlag *int
	BFlag *int
	CFlag *int
	cFlag *bool
	iFlag *bool
	vFlag *bool
	FFlag *bool
	nFlag *bool
)

// grep: grep [опции] шаблон [<путь к файлу или папке>]
// use: go run main.go -B 2 -A 4 -c -n -F -v -i 'auth' testdata.txt
// use 2: go run . -i downloads testdata.txt
// unnecessary comment
func main() {
	AFlag = flag.Int("A", 0, "Print NUM lines of trailing context after matching lines..")
	BFlag = flag.Int("B", 0, "Print NUM lines of leading context before matching lines.")
	CFlag = flag.Int("C", 0, "Print NUM lines of output context. Equivalent to '-A NUM -B NUM'.")
	cFlag = flag.Bool("c", false, "Suppress normal output; instead print a count of matching lines for each input file. With the '-v' option, count non-matching lines.") // количество строк, которые соответствуют шаблону
	iFlag = flag.Bool("i", false, "Ignore case distinctions in patterns and input data, so that characters that differ only in case match each other.")
	vFlag = flag.Bool("v", false, "Invert the sense of matching, to select non-matching lines.")
	FFlag = flag.Bool("F", false, "Interpret patterns as fixed strings, not regular expressions.")
	nFlag = flag.Bool("n", false, "Prefix each line of output with the 1-based line number within its input file.")

	if len(os.Args) < 3 {
		log.Fatal("error: pattern and input filenames are required")
		return
	}

	inputFileName := os.Args[len(os.Args)-1]
	pattern := os.Args[len(os.Args)-2]

	os.Args = os.Args[:len(os.Args)-2]

	flag.Parse()

	if *AFlag == 0 && *CFlag != 0 {
		*AFlag = *CFlag
	}

	if *BFlag == 0 && *CFlag != 0 {
		*BFlag = *CFlag
	}

	if len(os.Args) < flag.NFlag()+1 {
		log.Fatal("error: pattern and input filenames are required")
		return
	}

	matchingLinesCounter := 0

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal("error: can't open the file: ", err)
	}
	defer func(inputFile *os.File) {
		closeErr := inputFile.Close()
		if closeErr != nil {
			log.Fatal("error: can't close the file: ", closeErr)
		}
	}(inputFile)

	_ = strings.Trim(pattern, "'")

	if *iFlag {
		pattern = strings.ToLower(pattern)
	}

	var re *regexp.Regexp
	if !*FFlag {
		re = regexp.MustCompile(pattern)
	}

	type numberedLine struct {
		text   string
		number int
	}
	BSlice := make([]numberedLine, 0, *BFlag)
	toPrint := 0
	lineNumber := 0

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if *iFlag {
			line = strings.ToLower(line)
		}

		match := false
		if *FFlag {
			match = strings.Contains(line, pattern)
		} else {
			match = re.MatchString(line)
		}

		if *vFlag {
			match = !match
		}

		if toPrint != 0 {
			if match {
				toPrint = *AFlag + 1
			}

			matchingLinesCounter++
			if !*cFlag {
				if *nFlag {
					fmt.Printf("%d. %s\n", lineNumber, line)
				} else {
					fmt.Println(line)
				}
			}

			toPrint--
			continue
		}

		if match {
			if *AFlag != 0 {
				toPrint = *AFlag
			}

			if *BFlag != 0 {
				for i := 0; i < len(BSlice); i++ {
					matchingLinesCounter++

					if !*cFlag {
						if *nFlag {
							fmt.Printf("%d. %s\n", BSlice[i].number, BSlice[i].text)
						} else {
							fmt.Println(BSlice[i].text)
						}
					}
				}

				BSlice = make([]numberedLine, 0, *BFlag)
			}

			matchingLinesCounter++
			if !*cFlag {
				if *nFlag {
					fmt.Printf("%d. %s\n", lineNumber, line)
				} else {
					fmt.Println(line)
				}
			}
		} else {
			if *BFlag != 0 {
				if len(BSlice) == *BFlag {
					BSlice = BSlice[1:*BFlag:*BFlag]
				}
				BSlice = append(BSlice, numberedLine{text: line, number: lineNumber})
			}
		}
	}

	if *cFlag {
		fmt.Println(matchingLinesCounter)
	}
}
