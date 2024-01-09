package nipejs

import (
	"io"
	"os"
	"path/filepath"

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

func scanFolder(tmpfilename string, foldername string) (io.Reader, int) {
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
			if foldername[len(foldername)-1] == '/' {
				relativePath = foldername + relativePath
				relativePaths = append(relativePaths, relativePath)
			} else {
				relativePath = foldername + "/" + relativePath
				relativePaths = append(relativePaths, relativePath)
			}
		}
		return nil
	})

	log.Debug().Msgf("Files: %s", relativePaths)
	if err != nil {
		log.Error().Msgf("%s", err)
	}

	r := createTMPfile(tmpfilename, relativePaths)
	return r, len(relativePaths)
}

func ReadFiles(results chan Results, files chan string) {
	regexfile, _ := os.Open(*regexf)
	for file := range files {
		jsprefile, err := os.Open(file)
		if err != nil {
			log.Error().Msgf("Unable to open file: %s", file)
			wg.Done()
			continue
		}
		jsfile, _ := io.ReadAll(jsprefile)

		matchRegex(string(jsfile), file, results, regexfile)
		jsprefile.Close()
		wg.Done()
	}
}
