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
			"explicit not ok string",
			"bad-consent-not-ok-error",
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
			"consent v1 not ok does not allow purpose 1 (set cookie)",
			"BONciguONcjGKADACHENAOLS1rAAPABgAEAAIA3",
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
			"really long consent v2 ok",
			"CPAavcCPAavcCAGABCFRBKCsAP_AAH_AAAqIHFNf_X_fb3_j-_59_9t0eY1f9_7_v-0zjgeds-8Nyd_X_L8X5mM7vB36pq4KuR4Eu3LBAQdlHOHcTUmw6IkVqTPsbk2Mr7NKJ7PEinMbe2dYGH9_n9XT_ZKY79_____7__-_____7_f__-__3_vp9V---wOJAIMBAUAgAEMAAQIFCIQAAQhiQAAAABBCIBQJIAEqgAWVwEdoIEACAxAQgQAgBBQgwCAAQAAJKAgBACwQCAAiAQAAgAEAIAAEIAILACQEAAAEAJCAAiACECAgiAAg5DAgIgCCAFABAAAuJDACAMooASBAPGQGAAKAAqACGAEwALgAjgBlgDUAHZAPsA_ACMAFLAK2AbwBMQCbAFogLYAYEAw8BkQDOQGeAM-EQHwAVABWAC4AIYAZAAywBqADZAHYAPwAgABGAClgFPANYAdUA-QCGwEOgIvASIAmwBOwCkQFyAMCAYSAw8Bk4DOQGfCQAYADgBzgN_CQTgAEAALgAoACoAGQAOAAeABAACIAFQAMIAaABqADyAIYAigBMgCqAKwAWAAuABvADmAHoAQ0AiACJgEsAS4AmgBSgC3AGGAMgAZcA1ADVAGyAO8AewA-IB9gH6AQAAjABQQClgFPAL8AYoA1gBtADcAG8AOIAegA-QCGwEOgIqAReAkQBMQCZQE2AJ2AUOApEBYoC2AFyALvAYEAwYBhIDDQGHgMiAZIAycBlwDOQGfANIAadA1gDWQoAEAYQaBIACoAKwAXABDADIAGWANQAbIA7AB-AEAAIKARgApYBT4C0ALSAawA3gB1QD5AIbAQ6Ai8BIgCbAE7AKRAXIAwIBhIDDwGMAMnAZyAzwBnwcAEAA4Bv4qA2ABQAFQAQwAmABcAEcAMsAagA7AB-AEYAKXAWgBaQDeAJBATEAmwBTYC2AFyAMCAYeAyIBnIDPAGfANyHQWQAFwAUABUADIAHAAQAAiABdADAAMYAaABqADwAH0AQwBFACZAFUAVgAsABcADEAGYAN4AcwA9ACGAERAJYAmABNACjAFKALEAW4AwwBkADKAGiANQAbIA3wB3gD2gH2AfoBGACVAFBAKeAWKAtAC0gFzALyAX4AxQBuADiQHTAdQA9ACGwEOgIiAReAkEBIgCbAE7AKHAU0AqwBYsC2ALZAXAAuQBdoC7wGEgMNAYeAxIBjADHgGSAMnAZUAywBlwDOQGfANEgaQBpIDSwGnANYAbGPABAIqAb-QgZgALAAoABkAEQALgAYgBDACYAFUALgAYgAzABvAD0AI4AWIAygBqADfAHfAPsA_ACMAFBAKGAU-AtAC0gF-AMUAdQA9ACQQEiAJsAU0AsUBaMC2ALaAXAAuQBdoDDwGJAMiAZOAzkBngDPgGiANJAaWA4AlAyAAQAAsACgAGQAOAAigBgAGIAPAAiABMACqAFwAMQAZgA2gCGgEQARIAowBSgC3AGEAMoAaoA2QB3gD8AIwAU-AtAC0gGKANwAcQA6gCHQEXgJEATYAsUBbAC7QGHgMiAZOAywBnIDPAGfANIAawA4AmACARUA38pBBAAXABQAFQAMgAcABAACKAGAAYwA0ADUAHkAQwBFACYAFIAKoAWAAuABiADMAHMAQwAiABRgClAFiALcAZQA0QBqgDZAHfAPsA_ACMAFBAKGAVsAuYBeQDaAG4APQAh0BF4CRAE2AJ2AUOApoBWwCxQFsALgAXIAu0BhoDDwGMAMiAZIAycBlwDOQGeAM-gaQBpMDWANZAbGVABAA-Ab-A.YAAAAAAAAAAA", // nolinter: lll
		},
		{
			"bad consent v2 - wrong prefix, must start with C",
			"ONciguONcjGKADACHENAOCIAC0ta__AACiQABwAoABAACA",
		},
	}

	var all []string
	for _, c := range testcases {
		all = append(all, c.consent)
	}

	b.Run("all testcases", func(b *testing.B) {
		// on https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
		// section "A note on compiler optimisations"
		// we have a warning about the compiler may eliminate
		// ParseString function call, to prevent this we assign the result to
		// some variables out of the for loop scope
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
		b.Run("case "+tc.label, func(b *testing.B) {
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
		// on https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
		// section "A note on compiler optimisations"
		// we have a warning about the compiler may eliminate
		// ParseString function call, to prevent this we assign the result to
		// some variables out of the for loop scope
		var consent api.VendorConsents

		var err error
		for n := 0; n < b.N; n++ {
			consent, err = vendorconsent.ParseString(consents[n%max])
		}

		_ = consent
		_ = err
	})
}
