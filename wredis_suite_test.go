package wredis_test

import (
	"testing"

	. "github.com/crowdriff/wredis"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TestProcess is the root test process
func TestProcess(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "poolClient Suite")
}

// safe and unsafe are global pointer to poolClient
// objects used for testing
var (
	safe   Wredis
	unsafe Wredis
)

// BeforeSuite
var _ = BeforeSuite(func() {
	var err error

	// init safe poolClient
	safe, err = Safe()
	Ω(err).Should(BeNil())

	// init unsafe poolClient
	unsafe, err = Unsafe()
	Ω(err).Should(BeNil())
})

// AfterSuite
var _ = AfterSuite(func() {
	safe.Close()
	unsafe.Close()
})
