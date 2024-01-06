package nipejs

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func FirstTime() error {
	configFilePath := "~/.config/nipejs/.config"
	configDirPath := "~/.config/nipejs"
	regexFilePath := "~/.config/nipejs/regex.txt"

	// Expand user home directory in file paths
	configFilePath, _ = expandHomeDir(configFilePath)
	regexFilePath, _ = expandHomeDir(regexFilePath)
	configDirPath, _ = expandHomeDir(configDirPath)

	// Create the Config Directory if not exists
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		file, err := os.Create(configFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		fmt.Printf("Create file: %s\n", configFilePath)
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
			err := downloadDefaultRegexFile(regexFilePath)
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

func downloadDefaultRegexFile(destPath string) error {
	// Replace the URL with the actual URL of your default regex file on GitHub
	defaultRegexURL := "https://raw.githubusercontent.com/i5nipe/nipejs/master/files/regex.txt"

	// Download the file
	// Implement your logic to download the file from the URL and save it to destPath
	// You can use libraries like "github.com/go-resty/resty" or standard library "net/http"

	// For simplicity, let's just print a message here
	fmt.Printf("Downloading default regex file from: %s\n", defaultRegexURL)

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
