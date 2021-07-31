# JSScan

> Read list of JS files and look for sensitive data via regex. 

## ☕ Install
`go get github.com/i5nipe/nipejs`

## ☕ Usage examples

```
nipejs -urls jsfile -r regex.txt

nipejs -urls ~/Path/to/jsfile -s -r regex.txt

cat jsfile | nipejs -r regex.txt
```
