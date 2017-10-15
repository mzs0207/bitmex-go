package bitmex_test

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/joho/godotenv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/santacruz123/bitmex-go"
)

var _ = Describe("WebsocketPrivate", func() {
	var key, secret string

	log.SetLevel(log.DebugLevel)
	godotenv.Load()

	BeforeEach(func() {
		key = os.Getenv("BITMEX_KEY")
		if key == "" {
			Fail("Missing BITMEX_KEY variable")
		}

		secret = os.Getenv("BITMEX_SECRET")
		if secret == "" {
			Fail("Missing BITMEX_SECRET variable")
		}
	})

	AfterEach(func() {
		time.Sleep(time.Second)
	})

	It("Authenticate", func() {
		ws := bitmex.NewWS()
		err := ws.Connect()
		defer ws.Disconnect()

		Expect(err).Should(Succeed())
		_ = ws.Auth(key, secret)
		time.Sleep(time.Second)
	})

	It("Authenticate + chan", func() {
		ws := bitmex.NewWS()
		err := ws.Connect()
		defer ws.Disconnect()

		Expect(err).Should(Succeed())
		chAuth := ws.Auth(key, secret)

		select {
		case <-chAuth:
		case <-time.After(2 * time.Second):
			Fail("No auth signal received")
		}
	})

	It("Orders", func() {
		ws := bitmex.NewWS()
		err := ws.Connect()
		defer ws.Disconnect()

		Expect(err).Should(Succeed())

		chAuth := ws.Auth(key, secret)

		<-chAuth

		chOrder := make(chan bitmex.Order)
		_ = ws.SubOrder(chOrder, []bitmex.Contract{})

		select {
		case <-chOrder:
		case <-time.After(time.Second):
			Fail("No order received")
		}
	})

	It("Position", func() {
		ws := bitmex.NewWS()
		err := ws.Connect()
		defer ws.Disconnect()

		Expect(err).Should(Succeed())
		chAuth := ws.Auth(key, secret)

		<-chAuth

		chPosition := make(chan bitmex.WSPosition)
		_ = ws.SubPosition(chPosition, []bitmex.Contract{})

		select {
		case <-chPosition:
		case <-time.After(time.Second):
			Fail("No position received")
		}
	})
})
