# JSScan

> Read list of JS files and look for sensitive data via regex. 

You don't need to specify the regex file if you put it in `~/.nipe/regex.txt`

## ☕ Install
`go get github.com/i5nipe/nipejs`

## ☕ Usage examples

```
nipejs -urls jsfile -r regex.txt

nipejs -urls ~/Path/to/jsfile -s -r regex.txt

cat jsfile | nipejs -r regex.txt
```
