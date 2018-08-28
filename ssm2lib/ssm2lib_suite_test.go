package ssm2lib_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSsm2lib(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ssm2lib Suite")
}
