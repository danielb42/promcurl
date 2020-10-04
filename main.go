package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"regexp"

	wf "github.com/danielb42/whiteflag"
	c "github.com/logrusorgru/aurora/v3"
)

type (
	Line        interface{ Colorize() string }
	ScalarLine  struct{ Text string }
	VectorLine  struct{ Text string }
	HelpComment struct{ Text string }
	TypeComment struct{ Text string }
	OtherLine   struct{ Text string }
)

type Label struct {
	Key   string
	Value string
}

var (
	firstLine = true
	re_scalar = regexp.MustCompile(`^(\w+) ([0-9\-\.e\+Na]+)$`)
	re_vector = regexp.MustCompile(`^(\w+){(.*)} ([0-9\-\.e\+Na]+)$`)
	re_help   = regexp.MustCompile(`^(#) (HELP) (\w+) (.+)$`)
	re_type   = regexp.MustCompile(`^(#) (TYPE) (\w+) (\w+)$`)
	re_label  = regexp.MustCompile(`(\w+)="(.+?)"`)
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

	for s := bufio.NewScanner(resp.Body); s.Scan(); {
		line := parseLine(s.Text())
		fmt.Print(line.Colorize())
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

func parseLabels(line string) []Label {

	var labels []Label

	for _, matches := range re_label.FindAllStringSubmatch(line, -1) {
		label := Label{
			Key:   matches[1],
			Value: matches[2],
		}

		labels = append(labels, label)
	}

	return labels
}

func (line ScalarLine) Colorize() string {

	tokens := re_scalar.FindStringSubmatch(line.Text)

	return fmt.Sprintf("%s %s\n",
		c.Gray(20, tokens[1]),
		c.Green(tokens[2]),
	)
}

func (line VectorLine) Colorize() string {

	labels := ""
	for _, label := range parseLabels(line.Text) {
		labels += fmt.Sprintf("%s=%s ", c.Cyan(label.Key), c.BrightCyan(label.Value))
	}

	tokens := re_vector.FindStringSubmatch(line.Text)

	return fmt.Sprintf("%s %s%s\n",
		c.Gray(20, tokens[1]),
		labels,
		c.Green(tokens[3]),
	)
}

func (line HelpComment) Colorize() string {

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
		c.Gray(12, tokens[1]),
		c.Gray(12, tokens[2]),
		c.Gray(12, tokens[3]),
		c.Gray(20, tokens[4]),
	)
}

func (line TypeComment) Colorize() string {

	if wf.FlagPresent("nocomments") {
		return ""
	}

	tokens := re_type.FindStringSubmatch(line.Text)

	return fmt.Sprintf("%s %s %s %s\n",
		c.Gray(12, tokens[1]),
		c.Gray(12, tokens[2]),
		c.Gray(12, tokens[3]),
		c.Gray(20, tokens[4]),
	)
}

func (line OtherLine) Colorize() string {
	return fmt.Sprintf("%s\n", c.Red(line.Text))
}
