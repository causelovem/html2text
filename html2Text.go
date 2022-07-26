package html2Text

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	newLine   = "\n"
	spaceRune = ' '
)

// HTML2Text extracts text from html
func HTML2Text(htmlString string) string {
	var (
		// prevents from many new lines in a row
		canPrintNewline = false
		// prevents from many spaces in a row
		isSpaceNeeded = false
		// tells if was some space between tags
		wasSpace = false
		// unwanted tags counter
		skipTags = 0
	)

	// use tokenizer, not parser, because it faster, and we do not need html tree
	tokenizer := html.NewTokenizer(strings.NewReader(htmlString))
	clearString := strings.Builder{}
	clearString.Grow(len(htmlString))

	// writeString writes text to string builder
	writeString := func(text string) {
		if len(strings.TrimSpace(text)) > 0 {
			if isSpaceNeeded && wasSpace && text[0] != spaceRune {
				clearString.WriteRune(spaceRune)
				isSpaceNeeded = false
			}
			clearString.WriteString(text)
			canPrintNewline = true
			isSpaceNeeded = text[len(text)-1] != spaceRune
			wasSpace = false
		} else {
			wasSpace = true
		}
	}

	// writeNewLine writes new line without conditions, e.g. because of <br> tag
	writeNewLine := func() {
		if skipTags == 0 {
			clearString.WriteString(newLine)
			isSpaceNeeded = false
		}
	}

	// writeNewLineConditional writes new line only if needed
	writeNewLineConditional := func() {
		if skipTags == 0 && canPrintNewline {
			clearString.WriteString(newLine)
			canPrintNewline = false
			isSpaceNeeded = false
		}
	}

	// parse new token
	tokenType := tokenizer.Next()
	for tokenType != html.ErrorToken {
		switch tokenType {
		// if token is text - write it (skip empty strings)
		case html.TextToken:
			// do not move skipTags == 0 to writeString in order to avoid unnecessary tokenizer operations
			if skipTags == 0 {
				text := tokenizer.Token().Data
				writeString(text)
			}
		// add new line instead of some tags
		case html.StartTagToken:
			switch tokenizer.Token().DataAtom {
			case atom.Br, atom.Li:
				writeNewLine()
			case atom.P, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
				writeNewLineConditional()
			case atom.Noscript:
				tokenizer.Next()
				// because of bug in golang.org/x/net/html (all tokens inside <noscript> are TextToken)
				// we have to parse tags inside noscript tag one more time
				// do not move skipTags == 0 to writeString in order to avoid unnecessary recursion
				if skipTags == 0 {
					writeString(HTML2Text(tokenizer.Token().Data))
				}
			// we do not want to parse content from these tags, so skip them
			case atom.Head, atom.Script, atom.Style:
				skipTags++
			}
		// add new line instead of some tags
		case html.EndTagToken:
			switch tokenizer.Token().DataAtom {
			case atom.Ul:
				writeNewLine()
			case atom.P, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
				writeNewLineConditional()
			// end of unwanted tags
			case atom.Head, atom.Script, atom.Style:
				skipTags--
			}
		case html.SelfClosingTagToken:
			switch tokenizer.Token().DataAtom {
			case atom.Br, atom.Li:
				writeNewLine()
			}
		}

		// parse next token
		tokenType = tokenizer.Next()
	}

	return clearString.String()
}
