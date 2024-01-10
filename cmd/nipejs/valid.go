package nipejs

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	. "github.com/logrusorgru/aurora/v3"
)

func (resp Results) printresu() {
	fmt.Printf("\t%q\n", resp.Resu)
	defer wg.Done()
}

func (resp Results) printDefault(ident string) {
	fmt.Printf("\n%s %s %s%d%s\n", Cyan("[*]").Bold(),
		Magenta(resp.Url).Bold(), Cyan("["), resp.Len, Cyan("]"))

	fmt.Printf("%v\n", Cyan(fmt.Sprintf("Regex:  %s  %s", resp.Regex, Green(ident))))
}

func (resp Results) printb64() {
	sDec, _ := b64.StdEncoding.DecodeString(resp.Resu)
	resp.printresu()
	fmt.Printf("\t%s\n", Green(string(sDec)))
}

func (resp Results) printheruko() {
	req, err := http.NewRequest("POST", "https://api.heroku.com/apps", nil)
	if err != nil {
		resp.printresu()
		return
	}
	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resp.Resu))

	re, err := http.DefaultClient.Do(req)
	if err != nil {
		resp.printresu()
		return
	}
	defer re.Body.Close()

	html, err := io.ReadAll(re.Body)
	if err != nil {
		resp.printresu()
		return
	}

	resp.printresu()
	if string(html) != "{\"id\":\"unauthorized\",\"message\":\"Invalid credentials provided.\"}" {
		fmt.Printf("\t %s\n", Green(string(html)))
		// Deuboa(fmt.Sprintf("[NipeJS] Heroku API Key\nAPI:* %s *\nUrl:* %s *\nRegex:* %s *", resp.Resu, resp.Url, resp.Regex))
		// Deuboa(fmt.Sprintf("[NipeJS] Response Heroku\n\n%s", string(html)))
	} else {
		fmt.Printf("\t %s\n", Red(string(html)))
	}
}

func (resp Results) printmailgun() {
	req, err := http.NewRequest("GET", "https://api.mailgun.net/v3/domains", nil)
	if err != nil {
		resp.printresu()
		return
	}
	req.SetBasicAuth("api", resp.Resu)

	re, err := http.DefaultClient.Do(req)
	if err != nil {
		resp.printresu()
		return
	}
	defer re.Body.Close()

	html, err := io.ReadAll(re.Body)
	if err != nil {
		resp.printresu()
		return
	}

	resp.printresu()
	if string(html) != "{\"message\":\"Invalid private key\"}" {
		fmt.Printf("\t %s\n", Green(string(html)))
		// Deuboa(fmt.Sprintf("[NipeJS] API KEY Mailgun \nAPI:* %s *\nUrl:* %s *\nRegex:* %s *", resp.Resu, resp.Url, resp.Regex))
		// Deuboa(fmt.Sprintf("[NipeJS] Response Mailgun\n\n%s", string(html)))
	} else {
		fmt.Printf("\t %s\n", Red(string(html)))
	}
}

func (resp Results) printrecaptcha() {
	params := url.Values{}
	params.Add("secret", resp.Resu)
	params.Add("response", `test`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", body)
	if err != nil {
		resp.printresu()
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	re, err := http.DefaultClient.Do(req)
	if err != nil {
		resp.printresu()
		return
	}
	defer re.Body.Close()

	rvali := regexp.MustCompile(`invalid\-input\-secret`)
	html, err := io.ReadAll(re.Body)
	if err != nil {
		resp.printresu()
		return
	}
	baba := rvali.Match(html)

	if baba {
		resp.printresu()
		fmt.Printf("%s\n", Red(string(html)))
	} else {
		resp.printresu()
		fmt.Printf("%s\n", Green(string(html)))
		// Deuboa(fmt.Sprintf("[NipeJS] API Recaptcha\nAPI:* %s *\nUrl:* %s *\nRegex:* %s *", resp.Resu, resp.Url, resp.Regex))
		// Deuboa(fmt.Sprintf("[NipeJS] Response Recaptcha\n\n%s", string(html)))
	}
}
