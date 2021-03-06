package wredis_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	Context("FlushAll", func() {
		It("should not be able to FlushAll with a safe impl", func() {
			err := safe.FlushAll()
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: FlushAll requires unsafe impl. See wredis.Unsafe"))
		})

		It("should be able to FlushAll with an unsafe impl", func() {
			Ω(unsafe.FlushAll()).Should(Succeed())
		})
	})

	Context("FlushDB", func() {
		It("should not be able to FlushDb with a safe impl", func() {
			err := safe.FlushDB()
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: FlushDB requires unsafe impl. See wredis.Unsafe"))
		})

		It("should be able to FlushDb with an unsafe impl", func() {
			Ω(unsafe.FlushDB()).Should(Succeed())
		})
	})
})
