package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"prometheus-demo/chan/cybersource"
	"prometheus-demo/payment/card"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 1. define metrics
var (
	httpRequestsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and HTTP method and path.",
	}, []string{
		"code",
		"method",
		"path",
	})
	httpLatenciesHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_requests_latency",
		Help:    "How long (ms) HTTP requests processed, partitioned by status code and HTTP method and path.",
		Buckets: prometheus.LinearBuckets(100, 50, 10),
	}, []string{
		"code",
		"method",
		"path",
	})
)

var port = "8080"

func init() {
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	// 2. Metrics have to be registered to be exposed:
	prometheus.MustRegister(httpRequestsCounter)
	prometheus.MustRegister(httpLatenciesHistogram)
}

func main() {
	cardResource := card.Resource{
		Service: card.NewService(cybersource.NewMock()),
		// 4. use it
		// inject counter
		HttpRequestsCounter: httpRequestsCounter,
		// inject histogram
		HttpLatenciesHistogram: httpLatenciesHistogram,
	}
	restful.DefaultContainer.Add(cardResource.WebService())

	setupSwagger()
	setupMetrics()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// 3. expose url
func setupMetrics() {
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
}

func setupSwagger() {
	// generate swagger json file
	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// download swagger-ui/dist from https://github.com/swagger-api/swagger-ui/releases/tag/v3.41.1
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("swagger-ui/dist"))))

	// Optionally, you may need to enable CORS for the UI to work.
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)

	log.Printf("Get the API using http://localhost:%s/apidocs.json", port)
	log.Printf("Open Swagger UI using http://localhost:%s/apidocs/?url=http://localhost:%s/apidocs.json", port, port)
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "UserService",
			Description: "Resource for managing Users",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "john",
					Email: "john@doe.rp",
					URL:   "http://johndoe.org",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{{TagProps: spec.TagProps{
		Name:        "users",
		Description: "Managing users"}}}
}
