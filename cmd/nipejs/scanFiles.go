package nipejs

import (
	"os"
	"fmt"
	"bufio"
	"regexp"
	"io"

	log "github.com/projectdiscovery/gologger"
)

func ReadFiles(results chan Results,files chan string){
	rege, _ := getfile(*regexf)
	log.Debug().Msg("Abriu")

	for file := range files {
		jsprefile, err := os.Open(file)
		if err != nil {
			fmt.Println("Unable to open file:", *jsfilename)
			os.Exit(1)
		}
		jsfile, _ := io.ReadAll(jsprefile)

		scanner := bufio.NewScanner(rege)
		for scanner.Scan() {
			func(reges string) {
				log.Debug().Msg(scanner.Text())
				nurex := regexp.MustCompile(reges)
				bateu := nurex.FindString(string(jsfile))
				if bateu != "" {
					results <- Results{bateu, "FileName", reges, len(string(jsfile)) / 5}
				}
			}(scanner.Text())
		}
		wg.Done()
	}
	
}


