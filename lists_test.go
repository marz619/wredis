package wredis_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lists", func() {
	testList := "wredis::test::list"

	BeforeEach(func() {
		unsafe.Del(testList)
	})

	Context("LLen", func() {
		It("should return an error when no key provided", func() {
			_, err := safe.LLen("")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty key"))
		})

		It("should return 0 when the list doesn't exist", func() {
			i, err := safe.LLen(testList)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(i).Should(Equal(int64(0)))
		})

		It("should return the length of the list", func() {
			_, err := safe.LPush(testList, "1", "2", "3")
			Ω(err).ShouldNot(HaveOccurred())
			i, err := safe.LLen(testList)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(i).Should(Equal(int64(3)))
		})
	})

	Context("LPop", func() {
		It("should return an error when no key provided", func() {
			_, err := safe.LPop("")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty key"))
		})

		It("should return an error when popping from an empty list", func() {
			_, err := safe.LPop(testList)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("redigo: nil returned"))
		})

		It("should return the first item in a list", func() {
			n, err := safe.LPush(testList, "1", "2", "3")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(n).Should(Equal(int64(n)))

			i, err := safe.LPop(testList)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(i).Should(Equal("3"))
		})
	})

	Context("LPush", func() {
		It("should return an error when no key provided", func() {
			_, err := safe.LPush("")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty key"))
		})

		It("should return an error when no items provided", func() {
			_, err := safe.LPush(testList)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("must provide at least one item"))
		})

		It("should return an error when an item is empty", func() {
			_, err := safe.LPush(testList, "test", "")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("an item cannot be empty"))
		})

		It("should push an item to a new list", func() {
			n, err := safe.LPush(testList, "testing")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(n).Should(Equal(int64(1)))
		})

		It("should push multiple items to a new list", func() {
			n, err := safe.LPush(testList, "1", "2")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(n).Should(Equal(int64(2)))

			n, err = safe.LPush(testList, "3", "4")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(n).Should(Equal(int64(4)))
		})
	})

	Context("RPop", func() {
		It("should return an error when no key provided", func() {
			_, err := safe.RPop("")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty key"))
		})

		It("should return an error when popping from an empty list", func() {
			_, err := safe.RPop(testList)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("redigo: nil returned"))
		})

		It("should return the last item in a list", func() {
			n, err := safe.LPush(testList, "1", "2", "3")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(n).Should(Equal(int64(n)))

			i, err := safe.RPop(testList)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(i).Should(Equal("1"))
		})
	})

	Context("RPush", func() {
		It("should return an error when no key provided", func() {
			_, err := safe.RPush("")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty key"))
		})

		It("should return an error when no items provided", func() {
			_, err := safe.RPush(testList)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("must provide at least one item"))
		})

		It("should return an error when an item is empty", func() {
			_, err := safe.RPush(testList, "test", "")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("an item cannot be empty"))
		})

		It("should push an item to a new list", func() {
			n, err := safe.RPush(testList, "testing")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(n).Should(Equal(int64(1)))
		})

		It("should push multiple items to a new list", func() {
			n, err := safe.RPush(testList, "1", "2")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(n).Should(Equal(int64(2)))

			n, err = safe.RPush(testList, "3", "4")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(n).Should(Equal(int64(4)))
		})
	})
})
