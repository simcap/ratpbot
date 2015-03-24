package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

func BotReply(question string) *Reply {
	reply := new(Reply)

	reply.Question = question
	reply.Answer = detectors.run(question).Text()

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
	GlobalTrafficDetector = NewRegexpDetector(`^\s*(\?)\s*$`)
	MetroTrafficDetector  = NewRegexpDetector(`(?i)^\b*(metro)\b*\s*\??$`)
	RerTrafficDetector    = NewRegexpDetector(`(?i)^\b*(rer)\b*\s*\??$`)
	MetroLineDetector     = NewRegexpDetector(`\b([123456789]\b|1[01234])\b`)
	RerLineDetector       = NewRegexpDetector(`\b(?:[Rr][Ee][Rr])?([ABCDE]){1}\b`)
	NoActionDetector      = &NullDetector{}
)

var detectors DetectorChain = []Detector{
	MetroLineDetector,
	RerLineDetector,
	GlobalTrafficDetector,
	RerTrafficDetector,
	MetroTrafficDetector,
}

type answering func(text string) Answer

var Responses = responsesTable()

func responsesTable() map[Detector]answering {
	return map[Detector]answering{
		MetroLineDetector: func(text string) Answer {
			return SimpleAnswer{fmt.Sprintf("Ligne %s", text)}
		},
		RerLineDetector: func(text string) Answer {
			return SimpleAnswer{fmt.Sprintf("Rer %s", text)}
		},
		GlobalTrafficDetector: func(text string) Answer {
			message := GlobalTraffic()
			return SimpleAnswer{message.Text}
		},
		RerTrafficDetector: func(text string) Answer {
			message := RerTraffic()
			return SimpleAnswer{message.Text}
		},
		MetroTrafficDetector: func(text string) Answer {
			message := MetroTraffic()
			return SimpleAnswer{message.Text}
		},
		NoActionDetector: func(text string) Answer {
			return NewAnswer(UnsureAnswer, UsageAnswer)
		},
	}
}

type DetectorChain []Detector

func (c *DetectorChain) run(text string) Answer {
	var detector Detector = NoActionDetector
	var detected string

	for _, d := range *c {
		if grabbed, ok := d.grab(text); ok {
			detector = d
			detected = grabbed
			break
		}
	}

	return Responses[detector](detected)
}

type Detector interface {
	grab(text string) (string, bool)
}

type RegexpDetector struct {
	reg *regexp.Regexp
}

func NewRegexpDetector(exp string) *RegexpDetector {
	return &RegexpDetector{regexp.MustCompile(exp)}
}

func (d *RegexpDetector) grab(text string) (string, bool) {
	if item := d.reg.FindStringSubmatch(text); len(item) > 1 {
		return item[1], true
	}
	return "", false
}

type NullDetector struct{}

func (nd *NullDetector) grab(text string) (string, bool) {
	return "", true
}
