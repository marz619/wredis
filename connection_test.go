package wredis_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connection", func() {
	Context("Select", func() {
		It("should successfully select a different database", func() {
			_, err := safe.Select(1)
			Ω(err).Should(BeNil())

			_, err = safe.Select(0)
			Ω(err).Should(BeNil())
		})

		It("should fail if called on an already Selected instance", func() {
			w, err := safe.Select(1)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(w).ShouldNot(BeNil())

			_, err = w.Select(2)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: no select"))
		})
	})
})
