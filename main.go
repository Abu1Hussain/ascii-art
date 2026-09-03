package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 0 {
		return
	}

	if len(os.Args) < 2 {
		return
	}
	// case: hello OR Hi, ...
	// the user gives exactly one argument in the terminal
	if len(os.Args) == 2 {

		// send the user's input to the AsciiArt function
		art, err := AsciiArt(os.Args[1])
		// if there is an error, stop the program
		if err != nil {
			fmt.Println(err)
			return
		}

		// print the generated ASCII art in the terminal
		fmt.Print(art)

		// also write the generated ASCII art to output.txt
		err = os.WriteFile("output.txt", []byte(art), 0o644)

		// if the output file cannot be written, stop the program
		if err != nil {
			fmt.Println(err)
			return
		}

		return
	}

	// case 2: the user has written input/output file
	if len(os.Args) == 3 {

		// input file name
		inputFile := os.Args[1]

		// output file name
		outputFile := os.Args[2]

		// read all the data from the input file
		data, err := os.ReadFile(inputFile)
		// if the input file cannot be read, stop the program
		if err != nil {
			fmt.Println(err)
			return
		}

		// convert the data from []byte to string
		// and send it to the ASCII-art function
		art, err := AsciiArt(string(data))
		// if there is an error when converting
		if err != nil {
			fmt.Println(err)
			return
		}

		// write the changes into the output file
		err = os.WriteFile(outputFile, []byte(art), 0o644)

		// if the output file cannot be written, stop the program
		if err != nil {
			fmt.Println(err)
			return
		}

		// also print the generated ASCII art in the terminal
		fmt.Print(art)

		return
	}

	// case if the user has entered an invalid input
	fmt.Println("invalid input")
}

func AsciiArt(input string) (string, error) {
	// replace the literal \n with an actual newline
	// "Hello\nWorld" => "Hello
	// World"
	input = strings.ReplaceAll(input, `\n`, "\n")

	// handle the quote case:
	// '\!" #$%&'"'"'()*+,-./'
	// check for the starting and ending quote
	if len(input) >= 2 &&
		((strings.HasPrefix(input, "'") && strings.HasSuffix(input, "'")) ||
			(strings.HasPrefix(input, `"`) && strings.HasSuffix(input, `"`))) {
		input = input[1 : len(input)-1]
	}

	// handle:
	// '\!" #$%&'"'"'()*+,-./'
	// replacing '"'"' with '
	input = strings.ReplaceAll(input, `'"'"'`, "'")

	// replace \" with "
	input = strings.ReplaceAll(input, `\"`, `"`)

	// read from standard.txt
	fontBytes, err := os.ReadFile("standard.txt")
	// if standard.txt cannot be opened or read, return error
	if err != nil {
		return "", err
	}

	// handle Windows newlines
	fontContent := strings.ReplaceAll(string(fontBytes), "\r\n", "\n")

	// split each line in standard.txt
	fontLines := strings.Split(fontContent, "\n")

	// for the output
	var out strings.Builder

	// split the input into lines
	lines := strings.Split(input, "\n")

	// process each line
	for _, line := range lines {

		// handle empty lines
		if line == "" {
			// preserve the newline
			out.WriteByte('\n')
			continue
		}

		// process each row in the file
		// each character has a height of 8 rows
		for row := 0; row < 8; row++ {

			// convert the word into individual characters
			for j, char := range line {

				// check if the character is supported
				if char < ' ' || char > '~' {
					return "", fmt.Errorf("invalid input: unsupported character")
				}

				// find the character in standard.txt
				index := row + (int(char-' ') * 9) + 1

				// add it to the output
				out.WriteString(fontLines[index])

				// add space between characters
				if j != len([]rune(line))-1 {
					out.WriteByte(' ')
				}
			}

			// move to the next row
			out.WriteByte('\n')
		}
	}

	// convert strings.Builder into a string
	// nil means there was no error
	return out.String(), nil
}
