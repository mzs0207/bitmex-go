package bitmex_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBitmexGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BitmexGo Suite")
}
