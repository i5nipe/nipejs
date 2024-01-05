package nipejs

import (
	"github.com/valyala/fasthttp"
)

func GetBody(curl chan string, results chan Results, c *fasthttp.Client) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	for url := range curl {

		req.SetRequestURI(url)

		c.Do(req, resp)

		html := resp.Body()

		matchRegex(string(html), url, results)

		wg.Done()
	}
}
