package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	maxWidth = flag.Int("max.width", 70, "max width to allow comments, beyond which they are wrapped around")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("fixcomments: ")
	flag.Parse()

	for _, arg := range flag.Args() {
		fi, err := os.Lstat(arg)
		if err != nil {
			log.Printf("can't stat file %q: %v", arg, err)
			continue
		}
		content, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Printf("can't read file %q: %v", arg, err)
			continue
		}

		wrapped := bytes.NewBuffer(nil)
		changes, err := wrapComments(wrapped, bytes.NewReader(content))
		if err != nil {
			log.Printf("wrapping comments in %q: %v", arg, err)
		}

		err = ioutil.WriteFile(arg, wrapped.Bytes(), fi.Mode())
		if err != nil {
			log.Printf("can't replace content of file %q: %v", arg, err)
		}
		log.Printf("%d line changed: %q", changes, arg)
	}
}

func wrapComments(out io.Writer, orig io.Reader) (changes int, err error) {

	br := bufio.NewReader(orig)
	var (
		line       string
		result     string
		extraWords string
	)
	for {
		if len(extraWords) >= *maxWidth {
			line = extraWords
			extraWords = ""
		} else {
			p, err := br.ReadBytes('\n')
			switch err {
			default:
				return changes, err
			case io.EOF:
				if len(extraWords) == 0 {
					return changes, nil
				}
			case nil:
				// nothing to do
			}
			line = string(p)
		}

		result, extraWords = processLine(line, extraWords)
		if _, err := out.Write([]byte(result)); err != nil {
			return changes, err
		}
		if len(extraWords) != 0 {
			changes++
		}
	}
}

func processLine(line, extraWords string) (output, extra string) {
	if strings.HasPrefix(line, "//") {
		return processCommentLine(line, extraWords, *maxWidth)
	}
	if strings.Contains(line, "//") {
		return processCodeLine(line, extraWords, *maxWidth)
	}

	// line is empty or not a comment itself
	if len(extraWords) == 0 {
		return line, ""
	}

	// return a comment line, then the line itself
	return extraWords + "\n" + line, ""
}

func processCommentLine(line, extraWords string, maxWidth int) (output, extra string) {
	// remove the // part
	line = strings.TrimPrefix(line, "// ")
	line = strings.TrimPrefix(line, "//")
	if len(extraWords) > 0 {
		line = extraWords + " " + line
	} else {
		line = "// " + line
	}

	// wrapping the comments didn't exceed max?
	if len(line) < maxWidth {
		return line, ""
	}

	// need to return the <maxWidth line, and the words that were
	// chomped off

	// walk back until the `// ` part, looking for the first
	// whitespace
	for i := imin(len(line)-1, maxWidth); i > 2; i-- {
		// if not a rune, walk a bit more backward
		if !utf8.RuneStart(line[i]) {
			continue
		}
		r, _ := utf8.DecodeRuneInString(line[i:])
		// walk back til you find a space
		if !unicode.IsSpace(r) {
			continue
		}
		extra := strings.TrimSpace("// " + line[i+1:])
		return line[:i] + "\n", extra

	}
	// no whitespace in the entire string, just give up and let it
	// be too long...
	return line, ""
}

func processCodeLine(line, extraWords string, maxWidth int) (output, extra string) {

	leadingWhitespace := strings.IndexFunc(line, func(r rune) bool { return !unicode.IsSpace(r) })
	commentBegins := strings.Index(line, "//")

	if len(strings.TrimSpace(extraWords)) != 0 {
		extraLeadingWhitespace := strings.IndexFunc(extraWords, func(r rune) bool { return !unicode.IsSpace(r) })
		extraCommentBegins := strings.Index(extraWords, "//")

		// there is leading whitespace and the whitespace ends when a comment
		// begins
		if extraLeadingWhitespace >= 0 && extraCommentBegins == extraLeadingWhitespace {
			log.Printf("line=%q, extraWords=%q", string(line), string(extraWords))
		}
		if !strings.HasSuffix(extraWords, "\n") {
			extraWords = extraWords + "\n"
		}
		return extraWords, line
	}

	comment := line[commentBegins:]
	comment = strings.Trim(comment, "// ")
	comment = strings.Trim(comment, "//")
	comment = strings.TrimSpace(comment)
	output = line[:leadingWhitespace] + "// " + comment + "\n"

	extra = line[:leadingWhitespace] + strings.TrimSpace(line[leadingWhitespace:commentBegins])

	if len(comment) < maxWidth {
		return output, extra
	}

	// our comment itself is still too long

	// walk back until the `// ` part, looking for the first
	// whitespace
	for i := imin(len(output)-1, maxWidth); i > leadingWhitespace+2; i-- {
		// if not a rune, walk a bit more backward
		if !utf8.RuneStart(output[i]) {
			continue
		}
		r, _ := utf8.DecodeRuneInString(output[i:])
		// walk back til you find a space
		if !unicode.IsSpace(r) {
			continue
		}
		extraComment := strings.TrimSpace("// " + output[i+1:])
		return output[:i] + "\n", extra + extraComment
	}

	return output, extra
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
