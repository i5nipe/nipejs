package nipejs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/projectdiscovery/gologger"
)

func createTMPfile(filename string, strings2write []string) io.Reader {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal().Msg("Unable to create file in /tmp directory")
	}
	defer file.Close()

	for _, str := range strings2write {
		_, err := file.WriteString(str + "\n")
		if err != nil {
			log.Fatal().Msg("Unable to write in file on /tmp directory")
		}
	}
	tmpfile, _ := os.Open(filename)
	return tmpfile
}

func scanFolder(tmpfilename string, foldername string) io.Reader {
	var relativePaths []string

	err := filepath.Walk(foldername, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path != foldername {
			relativePath, err := filepath.Rel(foldername, path)
			if err != nil {
				return err
			}
			relativePath = foldername + "/" + relativePath
			relativePaths = append(relativePaths, relativePath)
		}
		return nil
	})

	log.Debug().Msg(fmt.Sprintf("Files: %s", relativePaths))
	if err != nil {
		log.Fatal().Msg("Unable to read all files in the directory")
	}

	r := createTMPfile(tmpfilename, relativePaths)
	return r
}

func ReadFiles(results chan Results, files chan string) {
	rege, _ := getfile(*regexf)

	for file := range files {
		jsprefile, err := os.Open(file)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Unable to open file: %s", file))
			wg.Done()
			continue
		}
		jsfile, _ := io.ReadAll(jsprefile)

		scanner := bufio.NewScanner(rege)
		for scanner.Scan() {
			func(reges string) {
				log.Debug().Msg(scanner.Text())
				nurex := regexp.MustCompile(reges)
				matches := nurex.FindAllString(string(jsfile), -1)
				for _, match := range matches {
					results <- Results{match, file, reges, len(string(jsfile)) / 5}
				}
			}(scanner.Text())
		}
		wg.Done()
	}
}
