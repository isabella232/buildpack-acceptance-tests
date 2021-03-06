package integration_test

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App that uses Kafka", func() {
	var app *cutlass.App
	AfterEach(func() { app = DestroyApp(app) })

	Context("deploying a basic PHP app using RdKafka module", func() {
		Context("after the RdKafka module has been loaded into PHP", func() {
			BeforeEach(func() {
				app = cutlass.New(filepath.Join(testdata, "with_rdkafka"))
				app.SetEnv("COMPOSER_GITHUB_OAUTH_TOKEN", os.Getenv("COMPOSER_GITHUB_OAUTH_TOKEN"))

				logLevel, found := os.LookupEnv("LOG_LEVEL")
				app.SetEnv("BP_DEBUG", strconv.FormatBool(found))
				app.SetEnv("LOG_LEVEL", logLevel)

				PushAppAndConfirm(app)
			})

			// missing a composer.lock file which is now required by Composer
			//  also missing a WEBDIR, a WEBDIR will need to be added to get this working too
			It("logs that Producer could not connect to a Kafka server", func() {
				Expect(app.GetBody("/producer.php")).To(ContainSubstring("Kafka error: Local: Broker transport failure"))
			})
		})
	})
})
