package linters_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLinters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Linters Suite")
}
