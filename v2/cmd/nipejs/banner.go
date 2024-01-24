package nipejs

import (
	. "github.com/logrusorgru/aurora/v4"
	"github.com/projectdiscovery/gologger"
)

const banner = ` ███▄    █  ██▓ ██▓███  ▓█████ ▄▄▄██▀▀▀██████
 ██ ▀█   █ ▓██▒▓██░  ██▒▓█   ▀   ▒██ ▒██    ▒
▓██  ▀█ ██▒▒██▒▓██░ ██▓▒▒███     ░██ ░ ▓██▄
▓██▒  ▐▌██▒░██░▒██▄█▓▒ ▒▒▓█  ▄▓██▄██▓  ▒   ██▒
▒██░   ▓██░░██░▒██▒ ░  ░░▒████▒▓███▒ ▒██████▒▒
░ ▒░   ▒ ▒ ░▓  ▒▓▒░ ░  ░░░ ▒░ ░▒▓▒▒░ ▒ ▒▓▒ ▒ ░
░ ░░   ░ ▒░ ▒ ░░▒ ░      ░ ░  ░▒ ░▒░ ░ ░▒  ░ ░
   ░   ░ ░  ▒ ░░░          ░   ░ ░ ░ ░  ░  ░
         ░  ░              ░  ░░   ░       ░
`

// Version is the current version
const Version = `v2.0.1`

// showBanner is used to show the banner to the user
func Banner() {
	gologger.Print().Msgf("%s", Magenta(banner))
}
