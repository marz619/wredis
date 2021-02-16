package wredis_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	Context("FlushAll", func() {
		It("should not be able to FlushAll with a safe poolClient", func() {
			err := safe.FlushAll()
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: FlushAll requires unsafe poolClient. See wredis.Unsafe"))
		})

		It("should be able to FlushAll with an unsafe poolClient", func() {
			Ω(unsafe.FlushAll()).Should(Succeed())
		})
	})

	Context("FlushDB", func() {
		It("should not be able to FlushDb with a safe poolClient", func() {
			err := safe.FlushDB()
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: FlushDB requires unsafe poolClient. See wredis.Unsafe"))
		})

		It("should be able to FlushDb with an unsafe poolClient", func() {
			Ω(unsafe.FlushDB()).Should(Succeed())
		})
	})
})
