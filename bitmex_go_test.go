package bitmex_test

import (
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
})
