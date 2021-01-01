package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"path"
)

// This is a silly go application that will echo out the characters of a file in a well known location that matches up with the first argument. This is intended as a joke to help people learn how to program.
// It does this by using a shadow directory under "typewriter". Any files you try to "write" will use these files as a source. Banging on the keyboard at random will echo the file to the terminal and ding when done so you, the presenter know when people have "completed" the file ;)
func main()  {
	flag.Parse()
	fileLocation := flag.Arg(0)
	if len(fileLocation) == 0 {
		log.Fatal(`No arguments passed, but at least 1 output file is required.
To which file should I save the output?
tw [OUTPUT_FILE_PATH]`)
	}

	typewriterDir := "typewriter"
	sourceFile, err := os.Open(path.Join(typewriterDir, fileLocation))
	if err != nil {
		// no source file, tell the user
		log.Fatalf("oops! looks like you mistyped the output file: '%s'", fileLocation)
	}

	outputFileDir := path.Dir(fileLocation)
	err = os.MkdirAll(outputFileDir, 0755)
	if err != nil {
		log.Fatalf("unable to create folders for OUTPUT_FILE: %s", outputFileDir)
	}

	dstFile, err := os.OpenFile(fileLocation, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("oops! the file %s isn't writeable", fileLocation)
	}

	func() {
		terminalState, err := terminal.MakeRaw(0)
		if err != nil {
			log.Fatalln("setting stdin to raw: ", err)
		}
		defer func() {
			if err = terminal.Restore(0, terminalState); err != nil {
				log.Println("warning, failed to restore terminal:", err)
			}
		}()

		_, err = os.Stdout.WriteString("[Start typing below]\r\n")

		in := bufio.NewReader(os.Stdin)
		src := bufio.NewReader(sourceFile)
		for {
			var srcChar rune
			_, _, err = in.ReadRune()
			if err != nil {
				log.Println("unable to read from stdin")
			}
			srcChar, _, err = src.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("oops! couldn't read from the hidden file. Ignore me if you're learning")
			}
			if srcChar == '\n' {
				_, err = os.Stdout.WriteString("\r")
			}
			if srcChar != '\r' {
				_, err = os.Stdout.WriteString(string(srcChar))
				if err != nil {
					log.Fatalln("oops! couldn't echo to terminal")
				}
			}
			_, err = dstFile.WriteString(string(srcChar))
			if err != nil {
				log.Fatalln("oops! couldn't write to the output file")
			}
		}

		fmt.Print("\a")
		_, err = os.Stdout.WriteString("[Done!]\r\n")
	}()
	_ = os.Stdin.Close()
}
