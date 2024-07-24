package main

import (
	"os"
	"strings"

	"github.com/wincus/gitlocator"
)

func main() {

	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	gl, err := gitlocator.NewGitLocator(cwd)

	if err != nil {
		panic(err)
	}

	url, err := gl.GetURL()

	if err != nil {
		panic(err)
	}

	clean, err := gl.IsClean()

	if err != nil {
		panic(err)
	}

	var output strings.Builder

	output.WriteString(url)

	if clean {
		output.WriteString("\n")
	} else {
		output.WriteString(" (dirty)\n")
	}

	os.Stdout.WriteString(output.String())
}
