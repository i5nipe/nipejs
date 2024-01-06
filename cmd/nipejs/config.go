package nipejs

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func FirstTime() error {
	configFilePath := "~/.config/nipejs/.config"
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
		fmt.Printf("Create Dir: %s\n", configDirPath)
	}

	// Check if the configuration file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("It looks like this is the first time you're running Nipejs.")
		fmt.Println("Would you like to download the default regex file from GitHub? (Y/n)")

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')

		if answer == "n\n" {
			fmt.Println("No action taken. You can manually set up your configuration later.")
			createConfig(configFilePath)
		} else {
			err := downloadRegex(regexFilePath, defaultRegexURL)
			if err != nil {
				return err
			}
			createConfig(configFilePath)
			fmt.Printf("Default regex file downloaded and saved to %s.\n", regexFilePath)
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
	req, err := http.NewRequest("Get", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Location", fmt.Sprintf("Nipejs %s", Version))

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

	fmt.Printf("Downloaded default regex file from: %s\n", url)
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
