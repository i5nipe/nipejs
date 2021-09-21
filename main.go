package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"sync"

	. "github.com/logrusorgru/aurora/v3"
	log "github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

var (
	regexf  = flag.String("r", "~/.nipe/regex.txt", "Regex file")
	usera   = flag.String("u", "Mozilla/5.0 (Windows NT 12.0; rv:88.0) Gecko/20100101 Firefox/88.0", "User-Agent")
	silent  = flag.Bool("s", false, "Silent Mode")
	threads = flag.Int("t", 70, "Number of threads")
	urls    = flag.String("urls", "", "List of URLs to scan")
	debug   = flag.Bool("b", false, "Debug mode (For developers)")
)
var wg sync.WaitGroup

type Results struct {
	Resu  string
	Url   string
	Regex string
	Len   int
}

func init() {
	flag.Parse()

	if *debug {
		log.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	}

	if !*silent {
		Banner()
	} else {
		log.DefaultLogger.SetMaxLevel(levels.LevelSilent)

	}
	if *regexf == "~/.nipe/regex.txt" {
		user, err := user.Current()
		if err != nil {
			log.Fatal().Msg(fmt.Sprintf("%s", err))
		}
		*regexf = fmt.Sprintf("%s/.nipe/regex.txt", user.HomeDir)
	}
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	file, _ := os.Open(*urls)

	results := make(chan Results, *threads)
	curl := make(chan string, *threads)

	for w := 1; w < *threads; w++ {
		go GetBody(curl, results)
	}

	scanner := bufio.NewScanner(file)
	if *urls == "" {
		scanner = bufio.NewScanner(os.Stdin)
	}

	go func() {
		for {
			resp := <-results
			switch resp.Regex {
			case `AAAA[A-Za-z0-9_-]{7}:[A-Za-z0-9_-]{140}`:
				resp.printdefault("Firebase")
				resp.printresu()
			case `sq0csp-[ 0-9A-Za-z\-_]{43}|sq0[a-z]{3}-[0-9A-Za-z\-_]{22,43}`:
				resp.printdefault("Square oauth secret")
				resp.printresu()
			case `sqOatp-[0-9A-Za-z\-_]{22}|EAAA[a-zA-Z0-9]{60}`:
				resp.printdefault("Square access token")
				resp.printresu()
			case `AC[a-zA-Z0-9_\-]{32}`:
				resp.printdefault("Twilio account SID")
				resp.printresu()
			case `AP[a-zA-Z0-9_\-]{32}`:
				resp.printdefault("Twilio APP SID")
				resp.printresu()
			case `[A-Za-z0-9]{125}`:
				resp.printdefault("Facebook")
				resp.printresu()
			case `s3\.amazonaws.com[/]+|[a-zA-Z0-9_-]*\.s3\.amazonaws.com`:
				resp.printdefault("S3 bucket")
				resp.printresu()
			case `\b(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}\b`:
				resp.printdefault("IPv4")
				resp.printresu()
			case `[a-f0-9]{32}`:
				resp.printdefault("MD5 hash")
				resp.printresu()
			case `6L[0-9A-Za-z-_]{38}|^6[0-9a-zA-Z_-]{39}`:
				resp.printdefault("Google Recaptcha")
				resp.printrecaptcha()
			case `key-[0-9a-zA-Z]{32}`:
				resp.printdefault("Mailgun")
				resp.printmailgun()
			case `[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}`:
				resp.printdefault("UUID")
				resp.printresu()
			case `(eyJ|YTo|Tzo|PD[89]|aHR0cHM6L|aHR0cDo|rO0)[a-zA-Z0-9+/]+={0,2}`:
				resp.printdefault("Base64")
				resp.printb64()
			case "":
				break
			default:
				resp.printdefault("")
				resp.printresu()
			}
		}
	}()
	for scanner.Scan() {
		curl <- scanner.Text()
		wg.Add(1)
	}

	wg.Wait()

	close(results)
	close(curl)
	defer file.Close()
}

func GetBody(curl chan string, results chan Results) {
	rege, falha := getregexfile()
	if falha {
		fmt.Println("Unable to open regexps file")
		return
	}
	for url := range curl {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			wg.Done()
			continue
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", *usera)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			wg.Done()
			continue
		}
		html, err := io.ReadAll(resp.Body)
		if err != nil {
			wg.Done()
			continue
		}

		scanner := bufio.NewScanner(rege)
		log.Debug().Msg(fmt.Sprintf("%v %s", Red("Url"), url))
		for scanner.Scan() {
			func(reges string) {
				log.Debug().Msg(scanner.Text())
				nurex := regexp.MustCompile(reges)
				bateu := nurex.FindString(string(html))
				if bateu != "" {
					results <- Results{bateu, url, reges, len(html) / 5}
				}
			}(scanner.Text())
		}
		resp.Body.Close()
		wg.Done()
	}
}

func getregexfile() (*os.File, bool) {
	rege, err := os.Open(*regexf)
	if err != nil {
		return rege, true
	}
	return rege, false
}
