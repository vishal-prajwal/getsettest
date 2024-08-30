package local_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bitbucket.org/junglee_games/getsetgo/pandora/cache/adapters/local"
)

var Adapter *Local

func TestLocal(t *testing.T) {
	InitializeAdapter()
	RegisterFailHandler(Fail)

	RunSpecs(t, "Local Suite")
}

func InitializeAdapter() {
	Adapter = Initialize(LocalAdapterConfig{})
}
