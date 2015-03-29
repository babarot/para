package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	//"golang.org/x/text/encoding"
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	//"golang.org/x/text/encoding/japanese"
)

const internal_default_style = "solarized"

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

func getStyle() (ret []string, err error) {
	python_code := `from pygments.styles import get_all_styles
s = list(get_all_styles())
print ' '.join(s)`

	c, err := checkPath("python")
	if err != nil {
		return
	}
	cmd := exec.Command(c)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	io.WriteString(stdin, python_code)
	stdin.Close()
	out, err := cmd.Output()
	if err != nil {
		return
	}
	ret = strings.Split(strings.TrimRight(string(out), "\n"), " ")

	return
}

func setStyle(s string) (ret string, err error) {
	styles, err := getStyle()
	if err != nil {
		return
	}

	if stringInSlice(s, styles) {
		ret = s
	} else {
		ret = "default"
	}

	return
}

func checkPath(cmd string) (ret string, err error) {
	ret, err = exec.LookPath(cmd)
	if err != nil {
		err = fmt.Errorf("%s: executable file not found in $PATH", cmd)
		return
	}

	return ret, nil
}

func main() {
	var style = flag.String("s", "", "pygmentize style")
	//flag.Usage = "a"
	flag.Parse()
	if *style == "" {
		*style = internal_default_style
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Input file is missing\n")
		os.Exit(1)
	}

	cmd, err := checkPath("pygmentize")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, item := range args {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()

			if !isExists(item) {
				fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", item)
				os.Exit(1)
			}

			style, err := setStyle(*style)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}

			out, err := exec.Command(
				cmd,
				"-O",
				"style="+style,
				"-f",
				"console256",
				"-g",
				item,
			).Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}

			//fmt.Print(outputData(out))
			msg, _ := eucjp_to_utf8(string(out))
			fmt.Print(msg)
		}(item)
	}
	wg.Wait()
}

func eucjp_to_utf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.EUCJP.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func sjis_to_utf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func utf8_to_sjis(str string) (string, error) {
	iostr := strings.NewReader(str)
	rio := transform.NewReader(iostr, japanese.ShiftJIS.NewEncoder())
	ret, err := ioutil.ReadAll(rio)
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func outputData(s []byte) (ret string) {
	ret = string(s)

	c, err := checkPath("nkf")
	if err != nil {
		return
	}
	cmd := exec.Command(c)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	io.WriteString(stdin, ret)
	stdin.Close()
	out, err := cmd.Output()
	if err != nil {
		return
	}
	ret = string(out)

	return
}

// にほんごであそぼ
// vim: fdm=marker
