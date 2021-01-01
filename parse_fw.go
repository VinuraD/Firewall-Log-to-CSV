package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func ParseLines(filePath string, parse func(string) (string, bool)) ([]string, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var results []string
	//var prev []string
	for scanner.Scan() { //%$#@^&!!!!
		if output, add := parse(scanner.Text()); add {
			results = append(results, output)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func filter(line []string, rulename string) (out []string) {
	//for _, l := range line {}
	headers := [5]string{"cs2", "dvchost", "dst", "dpt", "src"}
	_, available := Find(line, "cs2="+rulename)
	//fmt.Println(i)
	var row []string

	if available {
		//fmt.Println("ok")
		//row = append(row, line[i])
		for _, header := range headers {
			//
			for _, l := range line {
				item := strings.Split(l, "=")
				pos, available := Find(item, header)
				if available {
					row = append(row, item[pos+1])
					//fmt.Println(item[pos+1])

				}
			}
		}
	}
	return row
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func writetxt(filtline string) {
	f, err := os.OpenFile("text.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err, "yyy")
	}
	defer f.Close()
	//fmt.Println(filtline)

	if _, err := f.WriteString(filtline + "\n"); err != nil {
		log.Println(err, "xxx")
	}
	return
}

func createcsv() {
	csvfile, err := os.Create("fw_report.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	headers := []string{"RULE", "FW NAME", "DESTINATION IP", "DESTINATION PORT", "SOURCE IP", "COUNT"}
	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Write(headers)
	csvwriter.Flush()
	return
}

/*func createcsvdat(rows []string) {

	f, err := os.OpenFile("fw_report.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := csv.NewWriter(f)
	fmt.Println(rows)
	if len(rows) > 0 {
		w.Write(rows)
	} else {
	}
	w.Flush()
	return
}*/
func createcsvdat(rows []string) {

	f, err := os.OpenFile("fw_report.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := csv.NewWriter(f)
	for _, row := range rows {
		if len(row) > 0 {
			w.Write(strings.Split(row, " "))
		} else {
		}
	}

	w.Flush()
	return
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			//fmt.Println("dup")
			return true
		}
	}
	return false
}

func removedup(s []string) []string {
	var final []string
	for _, line := range s {
		line = line + " " + strconv.Itoa(getcount(s, line))
		if !(contains(final, line)) && len(line) > 5 {
			//fmt.Println(line)
			final = append(final, line)
		} else {
		}
	}
	return final
}

func getcount(s []string, line string) int {
	count := 0
	for _, l := range s {
		if line == l {
			count = count + 1
		} else {
		}
	}
	return count
}

func filterrulename() string {
	if len(os.Args) < 3 {
		fmt.Println("Usage: line_parser; provide the file name and the rule name separated by spaces. If the rule itself has spaces please fix that before processing")
		return ""
	} else if len(os.Args) == 3 {
		if contains(strings.Split(os.Args[2], ""), ",") {
			return strings.Split(os.Args[2], ",")[0]
		}
		return os.Args[2]
	} else {
		fmt.Println("Usage: line_parser; provide the file name and the rule name separated by spaces. If the rule itself has spaces please fix that before processing")
		return ""
	}
}

var prev []string

func main() {
	fmt.Println("Started")
	/*if len(os.Args) < 3 {
		fmt.Println("Usage: line_parser; provide the file name and the rule name separated by spaces. If the rule itself has spaces please fix that before processing")
		return
	}*/
	rulename := filterrulename()
	fmt.Println(rulename)

	createcsv()

	//lines, err := ParseLines(os.Args[1], func(s string) (string, bool) {
	//return s, true
	olines, err := ParseLines(os.Args[1], func(s string) (string, bool) {
		//if strings.HasPrefix(s, "cs1=") {
		//count := 0
		list := strings.Fields(s)
		row := filter(list, rulename)
		rowcat := strings.Join(row, " ")
		if len(rowcat) > 0 {
			//prev = append(prev, rowcat)
			//fmt.Println(prev)
			//list := strings.Fields(s)<<<<<<<<<<<
			//fmt.Println(len(list))
			//fmt.Println("**********************", list)
			//fmt.Println("found in Slice at position: ", sort.StringSlice(list).Search("cs1="))
			//row := filter(list, rulename)<<<<<<<<<<<
			//fmt.Println("*********", row)
			//writetxt(rowcat)
			return rowcat, true //s
		}
		return s, false
	})

	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}
	/*for _, l := range olines {
		fmt.Println(l)
		writetxt(l)
	}*/
	fmt.Println("Filtered from original log file")
	/*plines, err2 := ParseLines("text.txt", func(s string) (string, bool) {
		if true {
			return s, true //s
		}
		return s, false
	})
	if err2 != nil {
		fmt.Println("Error while parsing file", err)
		return
	}*/
	fmt.Println("Writing to CSV...")
	newlines := removedup(olines)
	sort.Strings(newlines)
	//fmt.Println(newlines)
	createcsvdat(newlines)
	os.Remove("text.txt")

	//for _, l := range lines {
	//	fmt.Println(l)
	//}
}
