package ratp

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Message struct {
	Text    string
	Timeout bool
	err     error
}

func RerTraffic() *Message {
	return fetchTraffic("http://www.ratp.fr/informer/trafic/trafic.php?cat=2")
}

func MetroTraffic() *Message {
	return fetchTraffic("http://www.ratp.fr/informer/trafic/trafic.php?cat=1")
}

func GlobalTraffic() *Message {
	return fetchTraffic("http://www.ratp.fr/informer/trafic/trafic.php")
}

func fetchTraffic(url string) *Message {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	text := findTrafficDivText(content)
	return &Message{text, false, nil}
}

var (
	traffic = regexp.MustCompile(`(?s)<div class="trafic">(.*?)</div>`)
	tags    = regexp.MustCompile(`<[/!]?[-\w]+>`)
)

func findTrafficDivText(content []byte) string {
	if found := traffic.FindSubmatch(content); found != nil {
		return retainHTMLTextOnly(string(found[1]))
	}
	return ""
}

func retainHTMLTextOnly(s string) string {
	return strings.Trim(tags.ReplaceAllLiteralString(s, ""), ". \n\r\t")
}
