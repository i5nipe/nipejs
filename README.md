# NipeJS

## ☕ Overview

NipeJS is a powerful and user-friendly tool designed to automate the process of detecting JavaScript leaks through precise regex pattern scanning. Whether you're securing web applications or conducting security assessments, NipeJS streamlines the identification of potential data leaks within JavaScript code.

<img src="./files/NipeJS.jpeg" alt="alt text" width="550"/>

## ☕ Key Features
- 💡 **Automated Leak Detection:** Efficiently scan large codebases for sensitive information.
- ⚡ **Concurrent Scanning:** Process multiple URLs or files simultaneously for faster results.
- 🔍 **Regex Pattern Validation:** Ensure pattern accuracy with built-in validation.
- 🌟 **Special Regexs for API Keys:** Automatically validate API keys for added convenience.
- 🔓 **Base64 Decryption Patterns:** Decrypt Base64-encoded strings to uncover hidden information.

## ☕ Installation
```bash
go install github.com/i5nipe/nipejs@latest
```

## ☕ Usage
NipeJS supports various input methods, making it suitable for diverse scenarios, from one-time scans to automated security workflows.

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
