# JSScan

> Read list of JS files and look for sensitive data via regex. 


## ☕ Install
```bash
go get github.com/i5nipe/nipejs
```

## ☕ Regular expressions
> Download the file "[files/regex.txt](https://github.com/i5nipe/nipejs/blob/master/files/regex.txt)"

- You don't need to specify the regex file if you put it in `~/.nipe/regex.txt`.
- This tool has some special regex, like decrypt base64 strings.
- Automatically test the authenticity of some API keys and notify for telegram if valid. (~~Not sure about the results yet.~~)
  - [Creating Telegram bot](https://core.telegram.org/bots#3-how-do-i-create-a-bot)


## ☕ Usage examples

```
nipejs -urls jsfile -r regex.txt

nipejs -urls ~/Path/to/jsfile -s -r regex.txt

cat jsfile | nipejs -r regex.txt
```

## Credits
---
- [KeyHacks](https://github.com/streaak/keyhacks)
- [JSScanner](https://github.com/0x240x23elu/JSScanner)
- [ProjectDiscovery/gologger](https://github.com/projectdiscovery/gologger)
- [odomojuli/RegExAPI](https://github.com/odomojuli/RegExAPI)
- [l4yton/RegHex](https://github.com/l4yton/RegHex)
