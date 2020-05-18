package vendorconsent_test

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/prebid/go-gdpr/api"
	"github.com/prebid/go-gdpr/vendorconsent"
)

func BenchmarkParse(b *testing.B) {
	testcases := []struct {
		label   string
		consent string
	}{
		{
			"empty string",
			"",
		},
		{
			"explicit nok string",
			"bad-consent-nok-error",
		},
		{
			"bad encoded string",
			"bad-encoded!!!string",
		},
		{
			"short consent v1 ok",
			"BONV8oqONXwgmADACHENAO7pqzAAppY", // good
		},
		{
			"long consent v1 ok",
			"BOzZIf5OzZIgJAHABBFRDG-AAAAvRr_7__7-_9_-_f__9uj3Or_v_f__32ccL59v_h_7v-_7fi_-1nV4u_1vft9yfk1-5ctDztp507iakivXmqdeb1v_nz3_9pxP78k89r7337Ew_v8_v-b7BCPN9Y3v-8KA", // nolinter: lll
		},
		{
			"consent v1 nok",
			"BONciguONcjGKADACHENAOLS1rAAPABgAEAAIA3", // nok
		},
		{
			"consent v1 maxvendorid=0",
			"BOzZB5dOzZB5dADABAENABAAAAAAAA",
		},
		{
			"short consent v2 ok",
			"COzSDo9OzSDo9B9AAAENAiCAALAAAAAAAAAACOQAQCOAAAAA", // v2
		},
		{
			"long consent v2 ok",
			"COzSDo9OzSDo9B9AAAENAiCAALAAAAAAAAAACOQAQCOAAAAA.IF5EX2S5OI2tho2YdF7BEYYwfJxyigMgShgQIsS8NwIeFbBoGPmAAHBG4JAQAGBAkkACBAQIsHGBcCQABgIgRiRCMQEGMjzNKBJBAggkbI0FACCVmnkHS3ZCY70-6u__bA", // nolinter: lll
		},
		{
			"bad consent v2",
			"ONciguONcjGKADACHENAOCIAC0ta__AACiQABwAoABAACA",
		},
	}

	var all []string
	for _, c := range testcases {
		all = append(all, c.consent)
	}
	b.Run("all testcases", func(b *testing.B) {
		var consent api.VendorConsents
		var err error
		max := len(all)
		for n := 0; n < b.N; n++ {
			consent, err = vendorconsent.ParseString(all[n%max])
		}
		_ = consent
		_ = err
	})

	for _, tc := range testcases {
		tc := tc
		b.Run(fmt.Sprintf("case %s", tc.label), func(b *testing.B) {
			var consent api.VendorConsents
			var err error
			for n := 0; n < b.N; n++ {
				consent, err = vendorconsent.ParseString(tc.consent)
			}
			_ = consent
			_ = err
		})
	}
}

var consentFile string

func init() {
	flag.StringVar(&consentFile, "consent-file", "", "ascii consent file")
}
func BenchmarkVerify(b *testing.B) {
	if consentFile == "" {
		b.SkipNow()
	}
	readFile, err := os.Open(consentFile)
	if err != nil {
		b.FailNow() // abort
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var consents []string

	for fileScanner.Scan() {
		consents = append(consents, fileScanner.Text())
	}

	max := len(consents)
	b.Run(fmt.Sprintf("testing just parsing %d consents/string", max), func(b *testing.B) {
		var consent api.VendorConsents
		var err error
		for n := 0; n < b.N; n++ {
			consent, err = vendorconsent.ParseString(consents[n%max])
		}
		_ = consent
		_ = err
	})
}
