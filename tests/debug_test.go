package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database Debug", func() {
	It("should have users in the database", func() {
		var users []User
		err := db.Find(&users).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(users).ToNot(BeEmpty())

		By("Checking users data")
		for _, u := range users {
			GinkgoWriter.Printf("ID: %d, Name: %s, Email: %s\n", u.ID, u.Name, u.Email)
		}
	})

	It("should have products in the database", func() {
		var products []Product
		err := db.Find(&products).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(products).ToNot(BeEmpty())

		By("Checking products data")
		for _, p := range products {
			GinkgoWriter.Printf("ID: %d, Code: %s, Name: %s, Price: %.2f\n", p.ID, p.Code, p.Name, p.Price)
		}
	})
})
