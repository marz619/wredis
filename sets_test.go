package wredis_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sets", func() {
	var (
		testKey  = "wredis::test::sets::one"
		otherKey = "wredis::test::sets::two"

		testSet  = []string{"a", "b", "c"}
		otherSet = []string{"a", "b", "d", "e"}
	)

	BeforeEach(func() {
		Ω(unsafe.SAdd(testKey, testSet...)).Should(BeEquivalentTo(3))
		Ω(unsafe.SAdd(otherKey, otherSet...)).Should(BeEquivalentTo(4))
	})

	AfterEach(func() {
		Ω(unsafe.FlushAll()).Should(Succeed())
	})

	Context("SAdd", func() {
		It("should Add members to an existing set successfully", func() {
			Ω(safe.SAdd(testKey, otherSet...)).Should(BeEquivalentTo(2))
		})

		It("should fail if an empty slice is passed to SAdd", func() {
			_, err := safe.SAdd(testKey, []string{}...)
			Ω(err.Error()).Should(Equal("wredis: no members"))
		})
	})

	Context("SCard", func() {
		It("should return the correct count of members in a set", func() {
			Ω(safe.SCard(testKey)).Should(BeEquivalentTo(3))
			Ω(safe.SCard(otherKey)).Should(BeEquivalentTo(4))
		})

		It("should fail given an empty key", func() {
			_, err := safe.SCard("")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty key"))
		})
	})

	Context("SDiffStore", func() {
		diffKey := "wredis::test::sets::diff"

		It("should successfully store the difference of two sets correctly", func() {
			diff, err := safe.SDiffStore(diffKey, otherKey, testKey)
			Ω(err).Should(BeNil())
			Ω(diff).Should(BeEquivalentTo(2))
		})

		It("should successfully store the difference of two sets correctly", func() {
			diff, err := safe.SDiffStore(diffKey, testKey, otherKey)
			Ω(err).Should(BeNil())
			Ω(diff).Should(BeEquivalentTo(1))
		})

		It("should fail if empty dest is passed", func() {
			_, err := safe.SDiffStore("")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty dest"))
		})

		It("should fail if no set keys are passed", func() {
			_, err := safe.SDiffStore(diffKey)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: no set keys"))
		})

		It("should fail if any empty set keys are passed", func() {
			_, err := safe.SDiffStore(diffKey, "key", "", "otherKey")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty set keys"))
		})
	})

	Context("SMembers", func() {
		It("should returns the members of a set successfully", func() {
			members, err := safe.SMembers(testKey)
			Ω(err).Should(BeNil())
			Ω(members).Should(HaveLen(3))
			Ω(members).Should(ConsistOf(testSet))
		})

		It("should return an error if key passed is empty", func() {
			_, err := safe.SMembers("")
			Ω(err).ShouldNot(BeNil())
			Ω(err.Error()).Should(Equal("wredis: empty key"))
		})
	})

	Context("SUnionStore", func() {
		unionKey := "wredis::test::sets::union"

		It("should successfully store the union of two sets correctly", func() {
			union, err := safe.SUnionStore(unionKey, otherKey, testKey)
			Ω(err).Should(BeNil())
			Ω(union).Should(BeEquivalentTo(5))
		})

		It("should successfully store the union of two sets correctly", func() {
			union, err := safe.SUnionStore(unionKey, testKey, otherKey)
			Ω(err).Should(BeNil())
			Ω(union).Should(BeEquivalentTo(5))
		})

		It("should fail if empty dest is passed", func() {
			_, err := safe.SUnionStore("")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty dest"))
		})

		It("should fail if no set keys are passed", func() {
			_, err := safe.SUnionStore(unionKey)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: no set keys"))
		})

		It("should fail if any empty set keys are passed", func() {
			_, err := safe.SUnionStore(unionKey, "key", "", "otherKey")
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal("wredis: empty set keys"))
		})
	})
})
