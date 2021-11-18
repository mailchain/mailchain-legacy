// Copyright 2021 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rfc2822

import (
	"mime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// mailParseDate parses an RFC 5322 date string. Its a copy from /usr/local/go/src/net/mail/message.go
// unexplained issue where some dates were not being parsed by the built binary.
// Current thought is it's related to the build container. Re-vist this when that changes or is updated.
func mailParseDate(date string) (time.Time, error) {
	dateLayoutsBuildOnce.Do(buildDateLayouts)
	// CR and LF must match and are tolerated anywhere in the date field.
	date = strings.ReplaceAll(date, "\r\n", "")
	if strings.Index(date, "\r") != -1 {
		return time.Time{}, errors.New("mail: header has a CR without LF")
	}
	// Re-using some addrParser methods which support obsolete text, i.e. non-printable ASCII
	p := addrParser{date, nil}
	p.skipSpace()

	// RFC 5322: zone = (FWS ( "+" / "-" ) 4DIGIT) / obs-zone
	// zone length is always 5 chars unless obsolete (obs-zone)
	if ind := strings.IndexAny(p.s, "+-"); ind != -1 && len(p.s) >= ind+5 {
		date = p.s[:ind+5]
		p.s = p.s[ind+5:]
	} else {
		ind := strings.Index(p.s, "T")
		if ind == 0 {
			// In this case we have the following date formats:
			// * Thu, 20 Nov 1997 09:55:06 MDT
			// * Thu, 20 Nov 1997 09:55:06 MDT (MDT)
			// * Thu, 20 Nov 1997 09:55:06 MDT (This comment)
			ind = strings.Index(p.s[1:], "T")
			if ind != -1 {
				ind++
			}
		}

		if ind != -1 && len(p.s) >= ind+5 {
			// The last letter T of the obsolete time zone is checked when no standard time zone is found.
			// If T is misplaced, the date to parse is garbage.
			date = p.s[:ind+1]
			p.s = p.s[ind+1:]
		}
	}
	if !p.skipCFWS() {
		return time.Time{}, errors.New("mail: misformatted parenthetical comment")
	}
	for _, layout := range dateLayouts {
		t, err := time.Parse(layout, date)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("mail: header could not be parsed")
}

type addrParser struct {
	s   string
	dec *mime.WordDecoder // may be nil
}

// skipCFWS skips CFWS as defined in RFC5322.
func (p *addrParser) skipCFWS() bool {
	p.skipSpace()

	for {
		if !p.consume('(') {
			break
		}

		if _, ok := p.consumeComment(); !ok {
			return false
		}

		p.skipSpace()
	}

	return true
}

// skipSpace skips the leading space and tab characters.
func (p *addrParser) skipSpace() {
	p.s = strings.TrimLeft(p.s, " \t")
}
func (p *addrParser) consume(c byte) bool {
	if p.empty() || p.peek() != c {
		return false
	}
	p.s = p.s[1:]
	return true
}
func (p *addrParser) consumeComment() (string, bool) {
	// '(' already consumed.
	depth := 1

	var comment string
	for {
		if p.empty() || depth == 0 {
			break
		}

		if p.peek() == '\\' && p.len() > 1 {
			p.s = p.s[1:]
		} else if p.peek() == '(' {
			depth++
		} else if p.peek() == ')' {
			depth--
		}
		if depth > 0 {
			comment += p.s[:1]
		}
		p.s = p.s[1:]
	}

	return comment, depth == 0
}

func (p *addrParser) peek() byte {
	return p.s[0]
}

func (p *addrParser) empty() bool {
	return p.len() == 0
}

func (p *addrParser) len() int {
	return len(p.s)
}

var (
	dateLayoutsBuildOnce sync.Once
	dateLayouts          []string
)

func buildDateLayouts() {
	// Generate layouts based on RFC 5322, section 3.3.

	dows := [...]string{"", "Mon, "}   // day-of-week
	days := [...]string{"2", "02"}     // day = 1*2DIGIT
	years := [...]string{"2006", "06"} // year = 4*DIGIT / 2*DIGIT
	seconds := [...]string{":05", ""}  // second
	// "-0700 (MST)" is not in RFC 5322, but is common.
	zones := [...]string{"-0700", "MST"} // zone = (("+" / "-") 4DIGIT) / "GMT" / ...

	for _, dow := range dows {
		for _, day := range days {
			for _, year := range years {
				for _, second := range seconds {
					for _, zone := range zones {
						s := dow + day + " Jan " + year + " 15:04" + second + " " + zone
						dateLayouts = append(dateLayouts, s)
					}
				}
			}
		}
	}
}
