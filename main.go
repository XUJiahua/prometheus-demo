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
)

var port = "8080"

func init() {
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}
}

func main() {
	cardResource := card.Resource{
		Service: card.NewService(cybersource.NewMock()),
	}
	restful.DefaultContainer.Add(cardResource.WebService())

	setupSwagger()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
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
