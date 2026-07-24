package tests

import (
	"fmt"
	"testing"
)

func TestDebugDatabase(t *testing.T) {
	// Check what's in the database
	var users []User
	db.Find(&users)
	fmt.Println("Users in database:")
	for _, u := range users {
		fmt.Printf("  ID: %d, Name: %s, Email: %s\n", u.ID, u.Name, u.Email)
	}

	var products []Product
	db.Find(&products)
	fmt.Println("Products in database:")
	for _, p := range products {
		fmt.Printf("  ID: %d, Code: %s, Name: %s\n", p.ID, p.Code, p.Name)
	}
}
