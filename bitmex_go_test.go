package bitmex_test

import (
	"os"
	"time"

	"github.com/apex/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/santacruz123/bitmex-go"
)

var _ = Describe("BitmexGo", func() {

	It("Trade", func() {
		SetDefaultEventuallyTimeout(time.Second)
		log.SetLevel(log.DebugLevel)

		ws := bitmex.NewWS()
		err := ws.Connect()
		Expect(err).Should(Succeed())

		ch := make(chan bitmex.WSTrade)
		ws.SubTrade(ch, []bitmex.Contracts{bitmex.XBTUSD})

		select {
		case <-ch:
		case <-time.After(time.Second):
			Fail("Nothing was received")
		}
	})

	It("Quote", func() {
		SetDefaultEventuallyTimeout(time.Second)
		log.SetLevel(log.DebugLevel)

		ws := bitmex.NewWS()
		err := ws.Connect()
		Expect(err).Should(Succeed())

		ch := make(chan bitmex.WSQuote)
		ws.SubQuote(ch, []bitmex.Contracts{bitmex.XBTUSD})

		select {
		case <-ch:
		case <-time.After(time.Second):
			Fail("Nothing was received")
		}
	})

	It("Authenticate", func() {
		log.SetLevel(log.DebugLevel)

		key, found := os.LookupEnv("BITMEX_KEY")

		if !found {
			Fail("Missing BITMEX_KEY variable")
		}

		secret, found := os.LookupEnv("BITMEX_SECRET")

		if !found {
			Fail("Missing BITMEX_SECRET variable")
		}

		ws := bitmex.NewWS()
		err := ws.Connect()
		Expect(err).Should(Succeed())
		_ = ws.Auth(key, secret)
		time.Sleep(time.Second)
	})

	It("Authenticate + chan", func() {
		log.SetLevel(log.DebugLevel)

		key, found := os.LookupEnv("BITMEX_KEY")

		if !found {
			Fail("Missing BITMEX_KEY variable")
		}

		secret, found := os.LookupEnv("BITMEX_SECRET")

		if !found {
			Fail("Missing BITMEX_SECRET variable")
		}

		ws := bitmex.NewWS()
		err := ws.Connect()
		Expect(err).Should(Succeed())
		chAuth := ws.Auth(key, secret)

		select {
		case <-chAuth:
		case <-time.After(2 * time.Second):
			Fail("No auth signal received")
		}
	})

	It("Orders", func() {
		log.SetLevel(log.DebugLevel)

		key, found := os.LookupEnv("BITMEX_KEY")

		if !found {
			Fail("Missing BITMEX_KEY variable")
		}

		secret, found := os.LookupEnv("BITMEX_SECRET")

		if !found {
			Fail("Missing BITMEX_SECRET variable")
		}

		ws := bitmex.NewWS()
		err := ws.Connect()
		Expect(err).Should(Succeed())
		chAuth := ws.Auth(key, secret)

		<-chAuth

		chOrder := make(chan bitmex.WSOrder)
		_ = ws.SubOrder(chOrder, []bitmex.Contracts{})

		select {
		case <-chOrder:
		case <-time.After(20 * time.Second):
			Fail("No order received")
		}
	})

	It("Position", func() {
		log.SetLevel(log.DebugLevel)

		key, found := os.LookupEnv("BITMEX_KEY")

		if !found {
			Fail("Missing BITMEX_KEY variable")
		}

		secret, found := os.LookupEnv("BITMEX_SECRET")

		if !found {
			Fail("Missing BITMEX_SECRET variable")
		}

		ws := bitmex.NewWS()
		err := ws.Connect()
		Expect(err).Should(Succeed())
		chAuth := ws.Auth(key, secret)

		<-chAuth

		chPosition := make(chan bitmex.WSPosition, 100)
		_ = ws.SubPosition(chPosition, []bitmex.Contracts{})

		select {
		case <-chPosition:
		case <-time.After(20 * time.Second):
			Fail("No position received")
		}
	})
})
