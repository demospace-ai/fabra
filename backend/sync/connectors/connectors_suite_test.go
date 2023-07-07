package connectors_test

import (
	"testing"

	"go.fabra.io/server/common/test"

	"gorm.io/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var db *gorm.DB
var cleanup func()

func TestConnectors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Connectors Suite")
}

var _ = BeforeSuite(func() {
	db, cleanup = test.SetupDatabase()
})

var _ = AfterSuite((func() {
	cleanup()
}))
