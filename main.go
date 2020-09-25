package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	wf "github.com/danielb42/whiteflag"
	c "github.com/logrusorgru/aurora/v3"
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	for s := bufio.NewScanner(resp.Body); s.Scan(); {
		line := s.Text()

		if !strings.HasPrefix(line, "#") {
			fmt.Fprintln(w, colorizeMetric(line))
			continue
		}

		if !wf.FlagPresent("nocomments") {
			if strings.HasPrefix(line, "# HELP") {
				fmt.Fprint(w, "\n", colorizeComment(line), "\n")
			} else if strings.HasPrefix(line, "# TYPE") {
				fmt.Fprintln(w, colorizeComment(line))
			}
		}
	}
}

func colorizeComment(text string) string {

	tokens := strings.Split(text, " ")

	if len(tokens) < 4 {
		return ""
	}

	return fmt.Sprintf("%s %s %s %s",
		c.Red("#").String(),
		c.Red(tokens[1]).String(),
		c.Yellow(tokens[2]).String(),
		c.Green(strings.Join(tokens[3:], " ")).String(),
	)
}

func colorizeMetric(text string) string {

	tokens := strings.Split(text, " ")

	if len(tokens) < 2 {
		return ""
	}

	return fmt.Sprintf("%s\t%s",
		c.Blue(tokens[0]).String(),
		c.White(tokens[1]).String(),
	)
}
