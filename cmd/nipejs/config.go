package nipejs

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	. "github.com/logrusorgru/aurora/v3"
	log "github.com/projectdiscovery/gologger"
)

func FirstTime() error {
	configFilePath := "~/.config/nipejs/.FirstRun"
	configDirPath := "~/.config/nipejs"
	regexFilePath := "~/.config/nipejs/regex.txt"
	defaultRegexURL := "https://raw.githubusercontent.com/i5nipe/nipejs/master/files/regex.txt"

	// Expand user home directory in file paths
	configFilePath, _ = expandHomeDir(configFilePath)
	regexFilePath, _ = expandHomeDir(regexFilePath)
	configDirPath, _ = expandHomeDir(configDirPath)

	// Create the Config Directory if not exists
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		err := os.MkdirAll(configDirPath, 0700)
		if err != nil {
			return err
		}
		log.Debug().Msgf("Create Dir: %s\n", configDirPath)
	}

	// Check if the configuration file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Info().
			Msgf("It looks like this is the %v you're running %v.", Magenta("first time").Bold(), Cyan("NipeJS").Bold())
		log.Info().
			Msgf("Would you like to %v the default regex file from GitHub? (%v/n)", Magenta("download").Bold(), Cyan("Y").Bold())

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')

		if answer == "n\n" {
			log.Debug().Msg("No action taken. You can manually set up your configuration later.")
			createConfig(configFilePath)
		} else {
			err := downloadRegex(regexFilePath, defaultRegexURL)
			if err != nil {
				return err
			}
			createConfig(configFilePath)
			log.Debug().Msgf("Default regex file downloaded and saved to %s.\n", regexFilePath)
		}
	}

	return nil
}

func expandHomeDir(path string) (string, error) {
	if path[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = filepath.Join(usr.HomeDir, path[2:])
	}
	return path, nil
}

func downloadRegex(destPath string, url string) error {
	// Create or open the destination file
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Download the file from the URL and save it to the destination file
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Location", fmt.Sprintf("https://i5nipe.com/nipejs/%s", Version))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Copy the content from the response to the destination file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Downloaded default regex file from: %s\n", url)
	return nil
}

func createConfig(configFilePath string) {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		file, err := os.Create(configFilePath)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
	}
}

func formatWithDots(value float64) string {
	// Convert the float to an integer
	intValue := int(value)

	// Convert the integer to a string
	strValue := strconv.Itoa(intValue)

	// Insert dots for every three digits from the end
	for i := len(strValue) - 3; i > 0; i -= 3 {
		strValue = strValue[:i] + "." + strValue[i:]
	}

	return strValue
}
