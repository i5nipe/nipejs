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
- Automatically test the authenticity of some API keys and notify the telegram if valid. (~~Not sure about the results yet.~~)

## ☕ Usage examples

```
nipejs -urls jsfile -r regex.txt

nipejs -urls ~/Path/to/jsfile -s -r regex.txt

cat jsfile | nipejs -r regex.txt
```
