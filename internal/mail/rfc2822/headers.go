package rfc2822

import (
	nm "net/mail"
	"time"

	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

func parseHeaders(h nm.Header) (*mail.Headers, error) {
	date, err := parseDate(h)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse `date`")
	}
	subject, err := parseSubject(h)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse `subject`")
	}
	to, err := parseTo(h)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse `to`")
	}
	from, err := parseFrom(h)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse `from`")
	}
	return &mail.Headers{
		Date:        *date,
		Subject:     subject,
		To:          *to,
		From:        *from,
		ContentType: parseContentType(h),
	}, nil
}

func parseTo(h nm.Header) (*mail.Address, error) {
	sources, ok := h["To"]
	if !ok {
		return nil, errors.Errorf("header missing")
	}
	if len(sources) == 0 {
		return nil, errors.Errorf("empty header")
	}

	return mail.ParseAddress(sources[0], "", "")
}
func parseFrom(h nm.Header) (*mail.Address, error) {
	sources, ok := h["From"]
	if !ok {
		return nil, errors.Errorf("header missing")
	}
	if len(sources) == 0 {
		return nil, errors.Errorf("empty header")
	}

	return mail.ParseAddress(sources[0], "", "")
}

func parseDate(h nm.Header) (*time.Time, error) {
	dateStrings, ok := h["Date"]
	if !ok {
		return nil, errors.Errorf("header missing")
	}
	if len(dateStrings) == 0 {
		return nil, errors.Errorf("empty header")
	}
	t, err := nm.ParseDate(dateStrings[0])
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseSubject(h nm.Header) (string, error) {
	sources, ok := h["Subject"]
	if !ok {
		return "", errors.Errorf("header missing")
	}
	if len(sources) == 0 {
		return "", errors.Errorf("empty header")
	}

	return sources[0], nil
}

func parseContentType(h nm.Header) string {
	sources, ok := h["Content-Type"]
	if !ok || len(sources) == 0 || sources[0] == "" {
		return mail.DefaultContentType
	}

	return sources[0]
}
