/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							go-carousel
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Command-line utility to verify compliance with UNIX-style
 * line endings (LF instead of CRLF).
 *-----------------------------------------------------------------*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

var flagOK bool

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/* checks file(s) for Unix line endings.
 */
func verifyUnixLineEndings(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("·Error opening file %s: %v\n", filePath, err)
		return false
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	hasCRLF := false
	var totalCRLF, totalLF, total int

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF { // @todo Catch Directories
			fmt.Printf("·Error reading file %s: %v\n", filePath, err)
			return false
		}

		if len(line) > 0 {
			// Check if the line ends with CRLF.
			if len(line) >= 2 && line[len(line)-2] == '\r' {
				totalCRLF++
				hasCRLF = true
			} else if len(line) >= 1 {
				totalLF++
			}
		}
		if err == io.EOF {
			break
		}
		total++
	}

	if hasCRLF {
		fmt.Printf("\t· '%s' (%d CRLF found).\n", filePath, totalCRLF)
		return false
	}

	if flagOK && (totalLF == total) {
		fmt.Printf("\t· '%s' is OK (LF Unix style)\n", filePath)
		return true
	}

	// Handles empty files or files without a final newline.
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		fmt.Printf("· '%s' is empty.\n", filePath)
		return true
	}

	//fmt.Printf("File '%s' is OK. No newline characters found.\n", filePath)
	return true
}

func Help() {
	fmt.Println("\t* * * goUnixStyle v1.0 * * *")
	fmt.Println("Usage:")
	fmt.Println("\tgoUnixStyle filename")
	fmt.Println("\tgoUnixStyle filename1 filename2 ...")
}

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/

func main() {
	var flagHelp bool
	flag.BoolVar(&flagHelp, "help", false, "Show help")
	flag.BoolVar(&flagOK, "ok", false, "Name files that are OK too")
	flag.Parse()

	if flagHelp {
		Help()
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		Help()
		os.Exit(1)
	}

	allClear := true
	fmt.Println("Checking for abnormal Windows CRLF endings...")
	for _, filePath := range flag.Args()[1:] {
		if !verifyUnixLineEndings(filePath) {
			allClear = false
		}
	}

	if allClear {
		fmt.Println("All specified files use only Unix line endings (LF).")
	} else {
		os.Exit(2)
	}
}
