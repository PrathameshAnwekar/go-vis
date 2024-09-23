package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FuncDecl struct {
	name       string
	parameters []string
	returnType []string
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please provide a project path.")
		return
	}
	fmt.Println("Parsing project:", args)

	projectRoot := args[1]
	goFileList, err := getGoFiles(projectRoot)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range goFileList {
		parseFile(file)
	}

	goFileList, err = removeCommonPrefix(goFileList)
	fmt.Println(goFileList)
}

func getGoFiles(projectRoot string) ([]string, error) {
	fmt.Println("Getting go files from:", projectRoot)
	fileList := make([]string, 0)
	err := filepath.WalkDir(projectRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") {
			fmt.Println("Adding: ", path)
			fileList = append(fileList, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileList, nil
}

func removeCommonPrefix(list []string) ([]string, error) {
	if len(list) <= 1 {
		return list, nil
	}
	prefix := list[0]
	for _, item := range list {
		for strings.Index(item, prefix) != 0 {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				break
			}
		}
	}
	for i, str := range list {
		list[i] = strings.TrimPrefix(str, prefix)
	}
	return list, nil
}

func parseFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	funcList := make([]FuncDecl, 0)
	for scanner.Scan() {
		word := ""
		for scanner.Scan() {
			if scanner.Text() != " " && scanner.Text() != "\n" {
				word = word + scanner.Text()
			} else {
				break
			}
		}
		switch word {
		case "func":
			funcList = append(funcList, parseFunc(scanner))
		default:
			// fmt.Println("Case not handled: ", word)
		}
		// fmt.Printf("%s;", word)
	}
	fmt.Println("LIST OF FUNCTIONS", funcList)
	return nil
}

func parseFunc(scanner *bufio.Scanner) FuncDecl {
	var funcName strings.Builder
	for scanner.Scan() {
		c := scanner.Text()
		if c != "(" {
			funcName.WriteString(c)
		} else {
			break
		}
	}
	inputParameterList := parseParameterList(scanner)
	returnValueList := parseReturnValues(scanner)

	return FuncDecl{name: funcName.String(), parameters: inputParameterList, returnType: returnValueList}
}

func parseParameterList(scanner *bufio.Scanner) []string {
	paranthesesCounter := 1
	parameter := ""
	parameterList := make([]string, 0)
	for scanner.Scan() {
		c := scanner.Text()
		if c == "," {
			if paranthesesCounter != 1 {
				parameter += c
			} else {
				parameterList = append(parameterList, parameter)
				parameter = ""
			}
			continue
		}
		if c == "(" {
			paranthesesCounter++
			parameter += c
			continue
		}
		if c == ")" {
			if paranthesesCounter == 1 {
				parameterList = append(parameterList, parameter)
				break
			} else {
				parameter += c
				paranthesesCounter--
			}
			continue
		}
		parameter += c
	}
	return parameterList
}

func parseReturnValues(scanner *bufio.Scanner) ([]string) {
	value := ""
	for scanner.Scan() {
		strings.TrimSpace(value)
		if scanner.Text() == "(" {
			return parseParameterList(scanner)
		}
		if scanner.Text() == "{" {
			return []string{value}
		}
		value += scanner.Text()
	}
	strings.TrimSpace(value)
	return []string{value}
}

func f() func(s string)() {
	return nil
}
