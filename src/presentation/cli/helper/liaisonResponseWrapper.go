package cliHelper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/goinfinite/os/src/presentation/liaison"
	"golang.org/x/term"
)

func LiaisonResponseWrapper(liaisonOutput liaison.LiaisonOutput) {
	exitCode := 0
	switch liaisonOutput.Status {
	case liaison.MultiStatus:
		exitCode = 1
	case liaison.UserError:
		exitCode = 1
	case liaison.InfraError:
		exitCode = 1
	}

	stdoutFileDescriptor := int(os.Stdout.Fd())
	isNonInteractiveSession := !term.IsTerminal(stdoutFileDescriptor)
	if isNonInteractiveSession {
		standardJsonBytes, err := json.Marshal(liaisonOutput)
		if err != nil {
			fmt.Println("ResponseEncodingError")
			os.Exit(1)
		}

		fmt.Println(string(standardJsonBytes))
		os.Exit(exitCode)
	}

	prettyJsonBytes, err := json.MarshalIndent(liaisonOutput, "", "  ")
	if err != nil {
		fmt.Println("PrettyResponseEncodingError")
		os.Exit(1)
	}

	syntaxHighlightingLexer := lexers.Get("json")
	if syntaxHighlightingLexer == nil {
		syntaxHighlightingLexer = lexers.Fallback
	}

	shIterator, err := syntaxHighlightingLexer.Tokenise(nil, string(prettyJsonBytes))
	if err != nil {
		fmt.Println("SyntaxHighlightingTokenizingError")
		os.Exit(1)
	}

	shFormatter := formatters.Get("terminal256")
	if shFormatter == nil {
		shFormatter = formatters.Fallback
	}

	err = shFormatter.Format(os.Stdout, styles.Vulcan, shIterator)
	if err != nil {
		fmt.Println("SyntaxHighlightingFormatError")
		os.Exit(1)
	}
	fmt.Println()
}
