package bot

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/simcap/ratpbot/ratp"
)

func Process(text string) *Reply {
	reply := new(Reply)

	reply.Question = text

	if item, ok := MetroLineDetector.item(text); ok {
		reply.Answer = SimpleAnswer{fmt.Sprintf("Ligne %s", item)}.Text()
	} else if item, ok := RerLineDetector.item(text); ok {
		reply.Answer = SimpleAnswer{fmt.Sprintf("Rer %s", item)}.Text()
	} else if _, ok := GlobalTrafficDetector.item(text); ok {
		message := ratp.GlobalTraffic()
		reply.Answer = SimpleAnswer{message.Text}.Text()
	} else if _, ok := RerTrafficDetector.item(text); ok {
		message := ratp.RerTraffic()
		reply.Answer = SimpleAnswer{message.Text}.Text()
	} else if _, ok := MetroTrafficDetector.item(text); ok {
		message := ratp.MetroTraffic()
		reply.Answer = SimpleAnswer{message.Text}.Text()
	} else {
		reply.Answer = NewAnswer(UnsureAnswer, UsageAnswer).Text()
	}

	return reply
}

type Reply struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Answer interface {
	Text() string
}

type SimpleAnswer struct {
	text string
}

func (a SimpleAnswer) Text() string {
	return a.text
}

type RandomAnswer struct {
	answers []Answer
}

func (a RandomAnswer) Text() string {
	index := rand.Intn(len(a.answers))
	return a.answers[index].Text()
}

var (
	UnsureAnswer = RandomAnswer{
		[]Answer{
			SimpleAnswer{"wtf?"},
			SimpleAnswer{"What do you mean?"},
			SimpleAnswer{"Hmm..."},
		},
	}

	UsageAnswer = SimpleAnswer{
		`Let me help you. Type '?' for global traffic info. For a line type the corresponding number/letter`,
	}
)

func NewAnswer(answers ...Answer) Answer {
	all := []string{}
	for _, a := range answers {
		all = append(all, a.Text())
	}
	return SimpleAnswer{strings.Join(all, ". ")}
}

var (
	GlobalTrafficDetector = Detector{
		regexp.MustCompile(`^\s*(\?)\s*$`),
	}
	MetroTrafficDetector = Detector{
		regexp.MustCompile(`(?i)^\b*(metro)\b*\s*\??$`),
	}
	RerTrafficDetector = Detector{
		regexp.MustCompile(`(?i)^\b*(rer)\b*\s*\??$`),
	}
	MetroLineDetector = Detector{
		regexp.MustCompile(`\b([123456789]\b|1[01234])\b`),
	}
	RerLineDetector = Detector{
		regexp.MustCompile(`\b(?:[Rr][Ee][Rr])?([ABCDE]){1}\b`),
	}
)

type Detector struct {
	reg *regexp.Regexp
}

func (d *Detector) item(text string) (string, bool) {
	if item := d.reg.FindStringSubmatch(text); len(item) > 1 {
		return item[1], true
	}
	return "", false
}
