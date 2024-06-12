package nipejs

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/user"
	"strings"
	"sync"
	"time"

	. "github.com/logrusorgru/aurora/v4"
	log "github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/valyala/fasthttp"
	"go.elara.ws/pcre"
)

var (
	regexf = flag.String("r", "~/.config/nipejs/regex.txt", "Regex file")
	usera  = flag.String(
		"a",
		"Mozilla/5.0 (Windows NT 12.0; rv:88.0) Gecko/20100101 Firefox/88.0",
		"User-Agent",
	)
	silent     = flag.Bool("s", false, "Silent Mode")
	threads    = flag.Int("c", 50, "Set the concurrency level")
	urls       = flag.String("u", "", "List of URLs to scan")
	debug      = flag.Bool("v", false, "Verbose mode")
	timeout    = flag.Int("timeout", 10, "Timeout in seconds")
	version    = flag.Bool("version", false, "Prints version information")
	jsdir      = flag.String("d", "", "Directory or File to match Regexs")
	Scan       = flag.Bool("no-scan", false, "Disable all scans for Special Regexs")
	jsonOutput = flag.Bool("json", false, "Enable json output")
)

var wg sync.WaitGroup

type Results struct {
	Match         string  `json:"Match"`
	Url           string  `json:"Url"`
	Regex         string  `json:"Regex"`
	Category      string  `json:"Category"`
	ContentLength float64 `json:"ContentLength"`
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
		Name: *usera,
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
			switch resp.Category {
			case "Google Recaptcha":
				resp.printSpecific("Google Recaptcha")

			case "Mailgun":
				resp.printSpecific("Mailgun")

			case "Base64":
				resp.printSpecific("Base64")

			case "empty":
				resp.printDefault("")

			case "":
				break

			default:
				resp.printDefault(resp.Category)
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

func matchRegex(target string, rlocation string, results chan Results, regexsfile []byte) {
	regexList := bufio.NewScanner(bytes.NewReader(regexsfile))
	for regexList.Scan() {
		lineText := regexList.Text()
		lineText = strings.TrimSpace(lineText)
		if lineText == "" {
			continue
		}
		parts := strings.Split(lineText, "\t\t")
		regex := parts[0]
		category := ""
		if len(parts) > 1 {
			category = strings.Join(parts[1:], "\t\t")
		}
		nurex := pcre.MustCompile(regex)

		matches := nurex.FindAllStringSubmatch(target, -1)
		for _, match := range matches {
			wg.Add(1)
      for _, unique_match := range match {

			category = strings.TrimSpace(category)
			if category == "" {
				category = "empty"
			}
			results <- Results{unique_match, rlocation, regex, category, float64(len(target)) / 1024}
      }
		}
	}
}

func calculateSeconds(startTimestamp, endTimestamp int64) float64 {
	startTime := time.Unix(0, startTimestamp)
	endTime := time.Unix(0, endTimestamp)

	duration := endTime.Sub(startTime)

	return duration.Seconds()
}

func checkRegexs(file string) {
	regexFile, err := os.Open(file)
	if err != nil {
		log.Fatal().Msgf("Unable to open regex file: %v", err)
	}
	defer regexFile.Close()

	regexCategories := make(map[string]string)
	regexL := bufio.NewScanner(regexFile)
	line := 1
	for regexL.Scan() {
		lineText := regexL.Text()
		lineText = strings.TrimSpace(lineText)
		if lineText == "" {
			continue
		}

		parts := strings.Split(lineText, "\t\t")
		regex := parts[0]
		category := ""
		if len(parts) > 1 {
			category = strings.Join(parts[1:], "\t\t")
		}

		_, err := pcre.Compile(regex)
		if err != nil {
			log.Fatal().
				Msgf("Regex on line %d not valid: %v", Cyan(line).Bold(), Red(lineText).Bold())
		}

		regexCategories[regex] = strings.TrimSpace(category)
		line++
	}

	if err := regexL.Err(); err != nil {
		log.Fatal().Msgf("Error reading regex file: %v", err)
	}
	for regexs, categories := range regexCategories {
		if categories == "" {
			log.Debug().Msgf("Regex: %v\n      Category: %v", Cyan(regexs), categories)
		} else {
			log.Debug().Msgf("Regex: %v\n      Category: %v", Cyan(regexs), Green(categories))
		}
	}
}
