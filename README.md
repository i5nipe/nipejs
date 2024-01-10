# NipeJS: Automated Regex Pattern Scanning for JavaScript Leaks

## Overview

NipeJS is a powerful and user-friendly tool designed to automate the process of detecting JavaScript leaks through precise regex pattern scanning. Whether you're securing web applications or conducting security assessments, NipeJS streamlines the identification of potential data leaks within JavaScript code.


> Read a list of JS files and look for sensitive data via regex.
<img src="./files/NipeJS.jpeg" alt="alt text" width="550"/>

## ☕ Installation
```bash
go install github.com/i5nipe/nipejs@latest
```

## ☕ Usage
NipeJS supports various input methods, including reading from standard input, scanning URLs from a file, or analyzing JavaScript files within a specified directory. The tool's flexibility makes it suitable for diverse scenarios, from one-time scans to automated security workflows.

## ☕ Example Commands

Scan URLs from STDIN: `cat UrlsList | nipejs`

Scan URLs from a file: `nipejs -u urlList.txt`

Analyze JavaScript files in a directory: `nipejs -d /path/to/js/files`

Analyze JavaScript file: `nipejs -d /path/to/js/files.js`

Specify a custom regex file: `nipejs -r regex.txt -d file.js`

## Contributing
Contributions to NipeJS are welcome! If you have suggestions, feature requests, or bug reports, please open an issue on GitHub.

## Acknowledgments

- [Elara](https://gitea.elara.ws/Elara6331/pcre)
- [KeyHacks](https://github.com/streaak/keyhacks)
- [JSScanner](https://github.com/0x240x23elu/JSScanner)
- [ProjectDiscovery/gologger](https://github.com/projectdiscovery/gologger)
- [odomojuli/RegExAPI](https://github.com/odomojuli/RegExAPI)
- [l4yton/RegHex](https://github.com/l4yton/RegHex)
- [logrusorgru/aurora](https://github.com/logrusorgru/aurora)
