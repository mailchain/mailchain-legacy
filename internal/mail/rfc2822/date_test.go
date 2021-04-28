package rfc2822

import (
	"net/mail"
	"testing"
	"time"
)

func TestDateParsing(t *testing.T) {
	tests := []struct {
		dateStr string
		exp     time.Time
	}{
		// RFC 5322, Appendix A.1.1
		{
			"Fri, 21 Nov 1997 09:55:06 -0600",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
		},
		// RFC 5322, Appendix A.6.2
		// Obsolete date.
		{
			"21 Nov 97 09:55:06 GMT",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("GMT", 0)),
		},
		// Commonly found format not specified by RFC 5322.
		{
			"Fri, 21 Nov 1997 09:55:06 -0600 (MDT)",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
		},
		{
			"Thu, 20 Nov 1997 09:55:06 -0600 (MDT)",
			time.Date(1997, 11, 20, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
		},
		{
			"Thu, 20 Nov 1997 09:55:06 GMT (GMT)",
			time.Date(1997, 11, 20, 9, 55, 6, 0, time.UTC),
		},
		{
			"Fri, 21 Nov 1997 09:55:06 +1300 (TOT)",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", +13*60*60)),
		},
	}
	for _, test := range tests {
		hdr := mail.Header{
			"Date": []string{test.dateStr},
		}
		date, err := hdr.Date()
		if err != nil {
			t.Errorf("Header(Date: %s).Date(): %v", test.dateStr, err)
		} else if !date.Equal(test.exp) {
			t.Errorf("Header(Date: %s).Date() = %+v, want %+v", test.dateStr, date, test.exp)
		}

		date, err = mailParseDate(test.dateStr)
		if err != nil {
			t.Errorf("mailParseDate(%s): %v", test.dateStr, err)
		} else if !date.Equal(test.exp) {
			t.Errorf("mailParseDate(%s) = %+v, want %+v", test.dateStr, date, test.exp)
		}
	}
}

func TestDateParsingCFWS(t *testing.T) {
	tests := []struct {
		dateStr string
		exp     time.Time
		valid   bool
	}{
		// FWS-only. No date.
		{
			"   ",
			// nil is not allowed
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			false,
		},
		// FWS is allowed before optional day of week.
		{
			"   Fri, 21 Nov 1997 09:55:06 -0600",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			true,
		},
		{
			"21 Nov 1997 09:55:06 -0600",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			true,
		},
		{
			"Fri 21 Nov 1997 09:55:06 -0600",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			false, // missing ,
		},
		// FWS is allowed before day of month but HTAB fails.
		{
			"Fri,        21 Nov 1997 09:55:06 -0600",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			true,
		},
		// FWS is allowed before and after year but HTAB fails.
		{
			"Fri, 21 Nov       1997     09:55:06 -0600",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			true,
		},
		// FWS is allowed before zone but HTAB is not handled. Obsolete timezone is handled.
		{
			"Fri, 21 Nov 1997 09:55:06           CST",
			time.Time{},
			true,
		},
		// FWS is allowed after date and a CRLF is already replaced.
		{
			"Fri, 21 Nov 1997 09:55:06           CST (no leading FWS and a trailing CRLF) \r\n",
			time.Time{},
			true,
		},
		// CFWS is a reduced set of US-ASCII where space and accentuated are obsolete. No error.
		{
			"Fri, 21    Nov 1997    09:55:06 -0600 (MDT and non-US-ASCII signs éèç )",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			true,
		},
		// CFWS is allowed after zone including a nested comment.
		// Trailing FWS is allowed.
		{
			"Fri, 21 Nov 1997 09:55:06 -0600    \r\n (thisisa(valid)cfws)   \t ",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			true,
		},
		// CRLF is incomplete and misplaced.
		{
			"Fri, 21 Nov 1997 \r 09:55:06 -0600    \r\n (thisisa(valid)cfws)   \t ",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			false,
		},
		// CRLF is complete but misplaced. No error is returned.
		{
			"Fri, 21 Nov 199\r\n7  09:55:06 -0600    \r\n (thisisa(valid)cfws)   \t ",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			true, // should be false in the strict interpretation of RFC 5322.
		},
		// Invalid ASCII in date.
		{
			"Fri, 21 Nov 1997 ù 09:55:06 -0600    \r\n (thisisa(valid)cfws)   \t ",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			false,
		},
		// CFWS chars () in date.
		{
			"Fri, 21 Nov () 1997 09:55:06 -0600    \r\n (thisisa(valid)cfws)   \t ",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			false,
		},
		// Timezone is invalid but T is found in comment.
		{
			"Fri, 21 Nov 1997 09:55:06 -060    \r\n (Thisisa(valid)cfws)   \t ",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			false,
		},
		// Date has no month.
		{
			"Fri, 21  1997 09:55:06 -0600",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			false,
		},
		// Invalid month : OCT iso Oct
		{
			"Fri, 21 OCT 1997 09:55:06 CST",
			time.Time{},
			false,
		},
		// A too short time zone.
		{
			"Fri, 21 Nov 1997 09:55:06 -060",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			false,
		},
		// A too short obsolete time zone.
		{
			"Fri, 21  1997 09:55:06 GT",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.FixedZone("", -6*60*60)),
			false,
		},
		// Ensure that the presence of "T" in the date
		// doesn't trip out ParseDate, as per issue 39260.
		{
			"Tue, 26 May 2020 14:04:40 GMT",
			time.Date(2020, 05, 26, 14, 04, 40, 0, time.UTC),
			true,
		},
		{
			"Tue, 26 May 2020 14:04:40 UT",
			time.Date(2020, 05, 26, 14, 04, 40, 0, time.UTC),
			false,
		},
		{
			"Thu, 21 May 2020 14:04:40 UT",
			time.Date(2020, 05, 21, 14, 04, 40, 0, time.UTC),
			false,
		},
		{
			"Thu, 21 May 2020 14:04:40 UTC",
			time.Date(2020, 05, 21, 14, 04, 40, 0, time.UTC),
			true,
		},
		{
			"Fri, 21 Nov 1997 09:55:06 GMT (GMT)",
			time.Date(1997, 11, 21, 9, 55, 6, 0, time.UTC),
			true,
		},
	}
	for _, test := range tests {
		hdr := mail.Header{
			"Date": []string{test.dateStr},
		}
		date, err := hdr.Date()
		if err != nil && test.valid {
			t.Errorf("Header(Date: %s).Date(): %v", test.dateStr, err)
		} else if err == nil && test.exp.IsZero() {
			// OK.  Used when exact result depends on the
			// system's local zoneinfo.
		} else if err == nil && !date.Equal(test.exp) && test.valid {
			t.Errorf("Header(Date: %s).Date() = %+v, want %+v", test.dateStr, date, test.exp)
		} else if err == nil && !test.valid { // an invalid expression was tested
			t.Errorf("Header(Date: %s).Date() did not return an error but %v", test.dateStr, date)
		}

		date, err = mailParseDate(test.dateStr)
		if err != nil && test.valid {
			t.Errorf("mailParseDate(%s): %v", test.dateStr, err)
		} else if err == nil && test.exp.IsZero() {
			// OK.  Used when exact result depends on the
			// system's local zoneinfo.
		} else if err == nil && !test.valid { // an invalid expression was tested
			t.Errorf("mailParseDate(%s) did not return an error but %v", test.dateStr, date)
		} else if err == nil && test.valid && !date.Equal(test.exp) {
			t.Errorf("mailParseDate(%s) = %+v, want %+v", test.dateStr, date, test.exp)
		}
	}
}
