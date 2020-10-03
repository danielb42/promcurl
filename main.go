package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	wf "github.com/danielb42/whiteflag"
	c "github.com/logrusorgru/aurora/v3"
)

type (
	Line        interface{ Format() string }
	ScalarLine  struct{ Text string }
	VectorLine  struct{ Text string }
	HelpComment struct{ Text string }
	TypeComment struct{ Text string }
	OtherLine   struct{ Text string }
)

var (
	firstLine = true
	re_scalar = regexp.MustCompile(`^([a-zA-Z0-9_]+) ([0-9\-\.e\+]+)$`)
	re_vector = regexp.MustCompile(`^([a-zA-Z0-9_]+){(.+)} ([0-9\-\.e\+]+)$`)
	re_help   = regexp.MustCompile(`^(#) (HELP) (\w+) (.+)$`)
	re_type   = regexp.MustCompile(`^(#) (TYPE) (\w+) (.+)$`)
)

func main() {

	wf.Alias("u", "url", "metrics endpoint url (https://myservice:9090/metrics)")
	wf.Alias("i", "insecure", "do not validate TLS certificate (ignore errors)")
	wf.Alias("n", "nocomments", "print only metric lines without comments")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: wf.FlagPresent("insecure"),
			},
		}}

	req, err := http.NewRequest("GET", wf.GetString("url"), nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	//TODO
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.Debug)
	//defer w.Flush()
	w := os.Stdout

	for s := bufio.NewScanner(resp.Body); s.Scan(); {
		line := parseLine(s.Text())
		fmt.Fprint(w, line.Format())
	}
}

func parseLine(line string) Line {
	if re_scalar.MatchString(line) {
		return ScalarLine{Text: line}
	}

	if re_vector.MatchString(line) {
		return VectorLine{Text: line}
	}

	if re_help.MatchString(line) {
		return HelpComment{Text: line}
	}

	if re_type.MatchString(line) {
		return TypeComment{Text: line}
	}

	return OtherLine{Text: line}
}

func (line ScalarLine) Format() string {

	tokens := re_scalar.FindStringSubmatch(line.Text)

	return fmt.Sprintf("%s %s\n",
		tokens[1],
		c.Green(tokens[2]).String(),
	)
}

func (line VectorLine) Format() string {

	tokens := re_vector.FindStringSubmatch(line.Text)

	return fmt.Sprintf("%s %s %s\n",
		tokens[1],
		c.Magenta(tokens[2]).String(),
		c.Green(tokens[3]).String(),
	)
}

func (line HelpComment) Format() string {

	if wf.FlagPresent("nocomments") {
		return ""
	}

	leadingNewline := ""

	if !firstLine {
		leadingNewline = "\n"
	} else {
		firstLine = false
	}

	tokens := re_help.FindStringSubmatch(line.Text)

	return fmt.Sprintf("%s%s %s %s %s\n",
		leadingNewline,
		c.Gray(12, tokens[1]).String(),
		c.Gray(12, tokens[2]).String(),
		c.Blue(tokens[3]).String(),
		c.Green(tokens[4]).String(),
	)
}

func (line TypeComment) Format() string {

	if wf.FlagPresent("nocomments") {
		return ""
	}

	tokens := re_type.FindStringSubmatch(line.Text)

	return fmt.Sprintf("%s %s %s %s\n",
		c.Gray(12, tokens[1]).String(),
		c.Gray(12, tokens[2]).String(),
		c.Blue(tokens[3]).String(),
		c.Green(tokens[4]).String(),
	)
}

func (line OtherLine) Format() string {
	return c.Red(line.Text).String() + "\n"
}
