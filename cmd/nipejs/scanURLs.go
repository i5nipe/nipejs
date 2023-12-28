package nipejs

import (
	"fmt"
	"bufio"
	"regexp"

	. "github.com/logrusorgru/aurora/v3"
	log "github.com/projectdiscovery/gologger"
	"github.com/valyala/fasthttp"
)

func GetBody(curl chan string, results chan Results, c *fasthttp.Client) {
  rege, _ := getfile(*regexf)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	for url := range curl {

		req.SetRequestURI(url)

		c.Do(req, resp)

		html := resp.Body()

		//var html need be a []byte
		scanner := bufio.NewScanner(rege)
		log.Debug().Msg(fmt.Sprintf("%v %s", Red("Url"), url))
		for scanner.Scan() {
			func(reges string) {
				log.Debug().Msg(scanner.Text())
				nurex := regexp.MustCompile(reges)
				bateu := nurex.FindString(string(html))
				if bateu != "" {
					results <- Results{bateu, url, reges, len(html) / 5}
				}
			}(scanner.Text())
		}
		wg.Done()
	}
}

