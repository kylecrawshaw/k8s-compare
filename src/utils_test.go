package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	Describe("contains function", func() {
		Context("when slice contains the item", func() {
			It("should return true", func() {
				slice := []string{"apple", "banana", "cherry"}
				result := contains(slice, "banana")
				Expect(result).To(BeTrue())
			})
		})

		Context("when slice does not contain the item", func() {
			It("should return false", func() {
				slice := []string{"apple", "banana", "cherry"}
				result := contains(slice, "orange")
				Expect(result).To(BeFalse())
			})
		})

		Context("when slice is empty", func() {
			It("should return false", func() {
				slice := []string{}
				result := contains(slice, "apple")
				Expect(result).To(BeFalse())
			})
		})

		Context("when item is empty string", func() {
			It("should return true if slice contains empty string", func() {
				slice := []string{"", "apple"}
				result := contains(slice, "")
				Expect(result).To(BeTrue())
			})

			It("should return false if slice does not contain empty string", func() {
				slice := []string{"apple", "banana"}
				result := contains(slice, "")
				Expect(result).To(BeFalse())
			})
		})
	})

	Describe("removeFromSlice function", func() {
		Context("when item exists in slice", func() {
			It("should remove the item", func() {
				slice := []string{"apple", "banana", "cherry"}
				result := removeFromSlice(slice, "banana")
				Expect(result).To(Equal([]string{"apple", "cherry"}))
			})

			It("should remove all occurrences", func() {
				slice := []string{"apple", "banana", "banana", "cherry"}
				result := removeFromSlice(slice, "banana")
				Expect(result).To(Equal([]string{"apple", "cherry"}))
			})
		})

		Context("when item does not exist in slice", func() {
			It("should return the original slice", func() {
				slice := []string{"apple", "banana", "cherry"}
				result := removeFromSlice(slice, "orange")
				Expect(result).To(Equal([]string{"apple", "banana", "cherry"}))
			})
		})

		Context("when slice is empty", func() {
			It("should return nil slice", func() {
				slice := []string{}
				result := removeFromSlice(slice, "apple")
				Expect(result).To(BeNil())
			})
		})

		Context("when removing empty string", func() {
			It("should remove empty string if it exists", func() {
				slice := []string{"", "apple", "banana"}
				result := removeFromSlice(slice, "")
				Expect(result).To(Equal([]string{"apple", "banana"}))
			})
		})
	})
})
