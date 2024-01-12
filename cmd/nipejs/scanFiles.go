package nipejs

import (
	"io"
	"io/ioutil"
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
	for file := range files {
		// Is not the best thing to too open the file every time for each file
		// But if is not in this way the default bufio.scanner from go don't run over the file again
		regexfile, _ := ioutil.ReadFile(*regexf)
		log.Debug().Msgf("Receveid file: %v", file)
		jsprefile, err := os.Open(file)
		if err != nil {
			log.Error().Msgf("Unable to open file: %s", file)
			wg.Done()
			continue
		}
		jsfile, _ := io.ReadAll(jsprefile)

		log.Debug().Msgf("File: %v\nContentLeagth: %v", file, len(jsfile))
		matchRegex(string(jsfile), file, results, regexfile)
		jsprefile.Close()
		wg.Done()
		log.Debug().Msgf("File Done: %v", file)
	}
}
