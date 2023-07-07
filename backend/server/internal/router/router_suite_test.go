package router_test

import (
	"testing"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/test"
	"go.fabra.io/server/internal/router"

	"gorm.io/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var db *gorm.DB
var r router.Router
var cleanup func()

func TestRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router Suite")
}

var _ = BeforeSuite(func() {
	db, cleanup = test.SetupDatabase()
	authService := auth.NewAuthService(db, test.MockCryptoService{})
	r = router.NewRouter(authService)
})

var _ = AfterSuite((func() {
	cleanup()
}))
