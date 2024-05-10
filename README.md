# NipeJS

## ☕ Overview

NipeJS is a powerful tool designed to detect JavaScript leaks through precise regex pattern scanning, streamlining the identification of potential data leaks within code.

<img src="./files/NipeJS.jpeg" alt="alt text" width="550"/>

## ☕ Key Features
- 💡 **Automated Leak Detection:** Efficiently scan large codebases for sensitive information.
- ⚡ **Concurrent Scanning:** Process multiple URLs or files simultaneously for faster results.
- 🌟 **Special Regexs for API Keys:** Automatically validate API keys for added convenience.
- 🔓 **Base64 Decryption Patterns:** Decrypt Base64-encoded strings to uncover hidden information.
- 🏷️ **Custom Regex Categories:** Dynamically categorize leaks by associating each regex with a custom category in the regex file.

## ☕ Additional Information
- **Editing Regexs File:** When adding a category to a regex, insert it after 2 tabs (`\t\t`). Be cautious, as some text editors may replace tabs with spaces.


## ☕ Installation
```bash
go install github.com/i5nipe/nipejs/v2@latest
```
The binary will be installed in the default Go binary directory. Ensure that this directory is included in your system's PATH variable to execute the nipejs command from any location in the terminal.

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## ☕ Example Commands

- Scan URLs from STDIN: `cat UrlsList | nipejs`
- Scan URLs from a file: `nipejs -u urlList.txt`
- Analyze JavaScript files in a directory: `nipejs -d /path/to/js/files`
- Analyze JavaScript file: `nipejs -d /path/to/js/files.js`
- Specify a custom regex file: `nipejs -r regex.txt -d file.js`

## ☕ Contributing
Contributions to NipeJS are welcome! If you have suggestions, feature requests, or bug reports, please [open an issue on GitHub](https://github.com/i5nipe/nipejs/issues).

## ☕ Acknowledgments

- [Elara](https://gitea.elara.ws/Elara6331/pcre)
- [KeyHacks](https://github.com/streaak/keyhacks)
- [JSScanner](https://github.com/0x240x23elu/JSScanner)
- [ProjectDiscovery/gologger](https://github.com/projectdiscovery/gologger)
- [odomojuli/RegExAPI](https://github.com/odomojuli/RegExAPI)
- [l4yton/RegHex](https://github.com/l4yton/RegHex)
- [logrusorgru/aurora](https://github.com/logrusorgru/aurora)
