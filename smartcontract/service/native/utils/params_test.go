package utils

import (
	"fmt"
	"testing"
)

func TestZptContractAddress(t *testing.T) {
	fromAddr := GovernanceContractAddress.ToBase58()
	fmt.Println(fromAddr)
}
