package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	//"log"
	"os"
	"strconv"
	"strings"
)

// Options - структура фалогов
type Options struct {
	C       bool
	D       bool
	U       bool
	I       bool
	F 		int
	S 		int

}

func execute(options *Options, src []string) []string {
	if options == nil || src == nil {
		return nil
	}
	text, stringsCount := gettext(options, src)
	text = checkOptions(options, text, stringsCount)
	return text
}

func gettext(options *Options, src []string) ([]string, map[string]int) {
	if len(src) == 0 {
		return nil, nil
	}
	var text []string
	stringsCount := make(map[string]int)

	prevString := src[0]
	prevTemplate := getStringTemplate(options, src[0])
	uniqCount := 1
	for _, s := range src[1:] {
		curTemplate := getStringTemplate(options, s)
		if prevTemplate == curTemplate {
			uniqCount++
		} else {
			text = append(text, prevString)
			stringsCount[prevString] = uniqCount
			prevString = s
			prevTemplate = curTemplate
			uniqCount = 1
		}
	}
	text = append(text, prevString)
	stringsCount[prevString] = uniqCount
	return text, stringsCount
}

func getStringTemplate(options *Options, s string) string {
	template := s

	if options.I {
		template = strings.ToLower(template)
	}

	if options.F != 0 {
		words := strings.Split(template, " ")
		if len(words) <= options.F {
			template = ""
		} else {
			template = strings.Join(words[options.F:], " ")
		}
	}

	if options.S != 0 {
		if len(template) <= options.S {
			template = ""
		} else {
			template = template[options.S:]
		}
	}
	return template
}

func checkOptions(options *Options, text []string, stringsCount map[string]int) []string {
	text = checkC(options, text, stringsCount)
	text = checkD(options, text, stringsCount)
	text = checkU(options, text, stringsCount)
	return text
}

func checkC(options *Options, text []string, stringsCount map[string]int) []string {
	var res []string
	if options.C {
		for _, s := range text {
			count := stringsCount[s]
			resString := strconv.Itoa(count)
			if s != "" {
				resString += " "
			}
			resString += s
			res = append(res, resString)
		}
		return res
	}
	return text
}

func checkD(options *Options, text []string, stringsCount map[string]int) []string {
	var res []string
	if options.D {
		for _, s := range text {
			count := stringsCount[s]
			if count > 1 {
				res = append(res, s)
			}
		}
		return res
	}
	return text
}

func checkU(options *Options, text []string, stringsCount map[string]int) []string {
	var res []string
	if options.U {
		for _, s := range text {
			count := stringsCount[s]
			if count == 1 {
				res = append(res, s)
			}
		}
		return res
	}
	return text
}



func initFlags(options *Options) error {
	flag.BoolVar(&options.C, "c", false, "for number of occurrences of lines in the input")
	flag.BoolVar(&options.D, "d", false, "print only those lines that were repeated in the input data")
	flag.BoolVar(&options.U, "u", false, "print only those lines that have not been repeated in the input data")
	flag.IntVar(&options.F, "f", 0, "ignore the first num_fields fields in the line")
	flag.IntVar(&options.S, "s", 0, "ignore the first num_chars characters in the string")
	flag.BoolVar(&options.I, "i", false, "case-insensitive")
	flag.Parse()

	if !options.C && options.D && options.U || options.C && !options.D && options.U || options.C && options.D && !options.U {
		return errors.New("invalid arguments passed")
	}
	return nil
}

func readFromFile(filename string) ([]string, error) {
	var inputFile io.Reader
	if filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		inputFile = file
	} else {
		inputFile = os.Stdin
	}

	var input []string
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return input, nil
}

func writeToFile(filename string, src []string) error {
	var outputFile io.Writer
	if filename != "" {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		outputFile = file
	} else {
		outputFile = os.Stdout
	}

	for _, s := range src {
		_, err := io.WriteString(outputFile, s+"\n")
		if err != nil {
			return err
		}
	}
	return nil
}


func main() {
	options := new(Options)
	err := initFlags(options)
	if err != nil {
		flag.PrintDefaults()
		os.Exit(2)
	}

	input, err := readFromFile(flag.Arg(0))
	if err != nil {
		os.Exit(1)
	}

	uniqRes := execute(options, input)

	if err := writeToFile(flag.Arg(1), uniqRes); err != nil {
		os.Exit(1)
	}
}