package wredis_test

import (
	. "github.com/crowdriff/wredis"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("poolClient", func() {
	var pool Wredis
	var err error

	AfterEach(func() {
		if pool != nil {
			Ω(pool.Close()).Should(Succeed())
		}
		pool = nil
	})

	It("should create a new default pool", func() {
		pool, err = Safe()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pool).ShouldNot(BeNil())
	})

	It("should fail to create a new pool given an invalid port", func() {
		_, err = Safe(Port(0))
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(Equal("wredis: invalid port"))
	})

	It("should create an unsafe pool successfully", func() {
		pool, err = Unsafe()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pool).ShouldNot(BeNil())
		Ω(pool.FlushAll()).Should(Succeed())
	})
})
