package nipejs

import (
	"bufio"
	"os"

	"github.com/valyala/fasthttp"
)

func GetBody(curl chan string, results chan Results, c *fasthttp.Client) {
	regexfile, _ := os.Open(*regexf)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	for url := range curl {

		req.SetRequestURI(url)

		c.Do(req, resp)

		html := resp.Body()

		matchRegex(string(html), url, results, regexfile)

		wg.Done()
	}
}

func countLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}
