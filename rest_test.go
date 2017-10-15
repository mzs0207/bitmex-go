package bitmex

import (
	"encoding/json"
	"fmt"
	"net/http/httputil"
	"os"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rest", func() {
	godotenv.Load()

	Context("Order", func() {
		It("Should sign", func() {
			// 'POST/api/v1/order1429631577995{"symbol":"XBTM15","price":219.0,"clOrdID":"mm_bitmex_1a/oemUeQ4CAJZgP3fjHsA","orderQty":98}'

			res := signature(
				"chNOOS4KvNXR_Xq4k4c9qsfoKWvnDecLATCRlcBwyKDYnWgO",
				"POST",
				"/order",
				1429631577995,
				[]byte(`{"symbol":"XBTM15","price":219.0,"clOrdID":"mm_bitmex_1a/oemUeQ4CAJZgP3fjHsA","orderQty":98}`),
			)

			Expect(res).To(Equal("93912e048daa5387759505a76c28d6e92c6a0d782504fc9980f4fb8adfc13e25"))
		})

		It("Should sign2", func() {
			// 'POST/api/v1/order1508108054{"symbol":"XBTUSD","orderQty":1}'

			res := signature(
				os.Getenv("BITMEX_SECRET"),
				"POST",
				"/order",
				1508108650,
				[]byte(`{"symbol":"XBTUSD","orderQty":1}`),
			)

			Expect(res).To(Equal("33746c29de52bfe80cd6e71dc3bae11faa31dc4f8acfd3f227b7e516d700e25d"))
		})

		It("Should make request object", func() {
			o := NewOrderMarket(XBTUSD, 1.9992232323)
			b := NewREST()

			body, err := json.Marshal(o)

			Expect(err).To(Succeed())

			req, err := b.request("POST", "/order", body)
			buf, err := httputil.DumpRequest(req, true)

			fmt.Println("Request start")
			fmt.Println(string(buf))
			fmt.Println("Request end")

			Expect(err).To(Succeed())
			Expect(buf).NotTo(BeEmpty())
		})

		It("Should send", func() {
			b := NewREST()

			o := NewOrderMarket(XBTUSD, 1)
			err := b.Send(o)

			Expect(err).To(Succeed())
		})
	})
})
