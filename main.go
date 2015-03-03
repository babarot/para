package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func isExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getStyle() string {
	python_code := `from pygments.styles import get_all_styles
s = list(get_all_styles())
print ' '.join(s)`

	cmd := exec.Command(checkPath("python"))
	stdin, err := cmd.StdinPipe()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	io.WriteString(stdin, python_code)
	stdin.Close()
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return strings.TrimRight(string(out), "\n")
}

func setStyle(s string) string {
	var style string

	style = "default"
	if stringInSlice(s, strings.Split(getStyle(), " ")) {
		style = s
	}

	return style
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func checkPath(cmd string) string {
	path, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "require: %s\n", cmd)
		os.Exit(1)
	}

	return path
}

func decisionStyle(internal string, external string) string {
	// No argumet is given
	// If the -s flag is not specified,
	// default string is assigned to the style variable
	if external != "default" {
		return setStyle(external)
	}

	return setStyle(internal)
}

func main() {
	var style = flag.String("s", "default", "pygmentize style")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Input file is missing\n")
		os.Exit(1)
	}

	for _, item := range args {
		if !isExists(item) {
			fmt.Fprintf(os.Stderr, item+": No such file or directory\n")
			os.Exit(1)
		}

		out, err := exec.Command(
			checkPath("pygmentize"),
			"-O",
			"style="+decisionStyle("solarized", *style),
			"-f",
			"console256",
			"-g",
			item,
		).Output()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	}
}
