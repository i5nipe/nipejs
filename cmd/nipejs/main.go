package nipejs

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/user"
	"regexp"
	"sync"
	"time"

	"github.com/dlclark/regexp2"
	. "github.com/logrusorgru/aurora/v3"
	log "github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/valyala/fasthttp"
)

var (
	regexf = flag.String("r", "~/.config/nipejs/regex.txt", "Regex file")
	usera  = flag.String(
		"a",
		"Mozilla/5.0 (Windows NT 12.0; rv:88.0) Gecko/20100101 Firefox/88.0",
		"User-Agent",
	)
	silent  = flag.Bool("s", false, "Silent Mode")
	threads = flag.Int("c", 50, "Set the concurrency level")
	urls    = flag.String("u", "", "List of URLs to scan")
	debug   = flag.Bool("b", false, "Debug mode")
	timeout = flag.Int("timeout", 10, "Timeout in seconds")
	version = flag.Bool("version", false, "Prints version information")
	jsdir   = flag.String("d", "", "Directory to scan all the files")
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
	if *version {
		fmt.Printf("NipeJS %s\n", Version)
		os.Exit(1)
	}

	if !*silent {
		Banner()
	} else {
		log.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	}
	if *regexf == "~/.config/nipejs/regex.txt" {
		user, err := user.Current()
		if err != nil {
			log.Fatal().Msgf("%s", err)
		}
		*regexf = fmt.Sprintf("%s/.config/nipejs/regex.txt", user.HomeDir)
	}
}

func Execute() {
	err := FirstTime()
	if err != nil {
		log.Error().Msgf("%v", err)
	}

	// Configs
	c := &fasthttp.Client{
		Name: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36",
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxConnWaitTimeout: time.Duration(*timeout) * time.Second,
	}
	urlsFile, _ := os.Open(*urls)

	checkRegexs(*regexf)

	allRegex, _ := countLines(*regexf)
	results := make(chan Results, *threads)
	curl := make(chan string, *threads)
	var input *bufio.Scanner
	var thread, countFiles, totalScan int
	StartTimestamp := time.Now().UnixNano()
	tmpFilename := fmt.Sprintf("/tmp/nipejs_%d%d", StartTimestamp, rand.Intn(100))

	// Switch case that define the Input Type
	switch {
	// If the input is STDIN (-u, -f or -d not especified)
	case *urls == "" && *jsdir == "":
		log.Debug().Msg("define input as Stdin")
		log.Debug().Msgf("Threads open: %d", *threads)
		for w := 0; w < *threads; w++ {
			go GetBody(curl, results, c)
		}
		input = bufio.NewScanner(os.Stdin)

	// If the input is for urls (-u especified)
	case *jsdir == "" && *urls != "":
		lines, _ := countLines(*urls)
		if lines < *threads {
			thread = lines
		} else {
			thread = *threads
		}
		log.Debug().Msgf("Threads open: %d", thread)
		for w := 0; w < thread; w++ {
			go GetBody(curl, results, c)
		}
		input = bufio.NewScanner(urlsFile)

		// If the input is for file or folder (-d)
	case *jsdir != "" && *urls == "":
		log.Debug().Msg(*jsdir)
		fileInfo, err := os.Stat(*jsdir)
		if err != nil {
			log.Fatal().Msg("Could not open Directory")
		}

		var tmpFile io.Reader

		// Scan a full Folder
		if fileInfo.IsDir() {
			tmpFile, countFiles = scanFolder(tmpFilename, *jsdir) // For directories
			if countFiles < *threads {
				thread = countFiles
			} else {
				thread = *threads
			}
			// Scan only one file
		} else {
			tmpFile = createTMPfile(tmpFilename, []string{*jsdir}) // For file
			thread = 1
		}
		defer os.Remove(tmpFilename)

		log.Debug().Msgf("Threads open: %d", thread)

		// Gouroutines That will wait for the input on channel 'curl'
		for w := 0; w < thread; w++ {
			go ReadFiles(results, curl)
		}

		input = bufio.NewScanner(tmpFile)

	default:
		log.Fatal().Msg("You can only specify one input method (-d or -u).")
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
			case `[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}`,
				`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`:
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
	for input.Scan() {
		// Send the input value to the functions that will match the regexs
		wg.Add(1)
		totalScan += 1
		curl <- input.Text()
	}
	wg.Wait()

	// Ending program
	close(results)
	close(curl)
	endTimestamp := time.Now().UnixNano()
	executionTime := calculateSeconds(StartTimestamp, endTimestamp)
	defer urlsFile.Close()
	fmt.Println("")
	log.Info().
		Msgf("Nipejs done: %d files with %d regex patterns scanned in %.2f seconds", Magenta(totalScan).Bold(), Cyan(allRegex).Bold(), Red(executionTime).Bold())
}

func matchRegex(target string, rlocation string, results chan Results, regexsfile io.Reader) {
	regexList := bufio.NewScanner(regexsfile)
	for regexList.Scan() {
		func(regex string) {
			nurex, err := regexp2.Compile(regex, 0)
			if err != nil {
				log.Fatal().Msg("Error in matchRegex")
			}
			matcher, err := nurex.FindStringMatch(target)
			if err != nil {
				log.Fatal().Msg("Error in matchRegex")
			}
			for matcher != nil {
				match := matcher.GroupByNumber(0).String()

				wg.Add(1)
				results <- Results{match, rlocation, regex, len(target) / 5}

				matcher, err = nurex.FindNextMatch(matcher)
				if err != nil {
					log.Fatal().Msg("Error Finding the Next Match")
				}
			}
		}(regexList.Text())
	}
}

func calculateSeconds(startTimestamp, endTimestamp int64) float64 {
	// Convert Unix nano timestamps to time.Time
	startTime := time.Unix(0, startTimestamp)
	endTime := time.Unix(0, endTimestamp)

	// Calculate the duration between two timestamps
	duration := endTime.Sub(startTime)

	return duration.Seconds()
}

func checkRegexs(file string) {
	regexFile, err := os.Open(file)
	if err != nil {
		log.Fatal().Msg("Unable to open regex file")
	}
	defer regexFile.Close()
	regexL := bufio.NewScanner(regexFile)
	line := 1
	for regexL.Scan() {
		_, err := regexp.Compile(regexL.Text())
		if err != nil {
			log.Fatal().
				Msgf("Regex on line %d not valid: %v", Cyan(line).Bold(), Red(regexL.Text()).Bold())
		}
		line += 1
	}
}
