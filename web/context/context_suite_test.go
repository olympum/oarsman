package context_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Context Suite")
}
