package e2e_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/MehmetTalhaSeker/mts-blog-api/e2e"
	postgresadapter "github.com/MehmetTalhaSeker/mts-blog-api/internal/adapter/postgres"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/database"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/testutils"
	"github.com/MehmetTalhaSeker/mts-blog-api/pkg/auth"
	"github.com/MehmetTalhaSeker/mts-blog-api/pkg/user"
)

var (
	container, p = testutils.NewPostgresTestContainer(context.Background())
	store        = database.NewPostgresStore(database.WithUser("test-user"), database.WithName("test-name"), database.WithPassword("test-password"), database.WithPort(p.Port()))
)

var e *echo.Echo

var _ = BeforeSuite(func() {
	done := make(chan struct{})

	go func() {
		rbac := rbac.New()

		store.InitDB()

		// initialize db repo.
		userRepo := postgresadapter.NewUserRepository(store.GetInstance())

		e = e2e.InitEcho()

		// create a new router group.
		routerGroup := e.Group("v1")

		// authentication router initialization.
		authRouter := &auth.Router{
			RouterGroup:    routerGroup,
			UserRepository: userRepo,
		}
		authRouter.New()

		// user router initialization.
		userRouter := &user.Router{
			Authenticate:   e2e.AuthMid(),
			RBAC:           rbac,
			RouterGroup:    routerGroup,
			UserRepository: userRepo,
		}
		userRouter.New()

		done <- struct{}{}
		err := e.Start(":8080")
		if err != nil {
			log.Printf("error while starting the server: %v\n", err)
		}
	}()

	<-done
	time.Sleep(1 * time.Second)
})

var _ = AfterSuite(func() {
	testutils.TerminateContainer(context.Background(), container)
	err := e.Close()
	if err != nil {
		log.Printf("error while closing the server: %v\n", err)
	}
})

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "E2E Suite")
}
