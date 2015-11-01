package model

import (
	"testing"
	"fmt"
)

func TestEntryID(t *testing.T) {
	entry := Entry{"Level", ""}
	fmt.Println("ID:", entry.ID())
	// Output ID: 2698725818
}