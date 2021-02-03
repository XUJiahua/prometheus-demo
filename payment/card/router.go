package card

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

type Resource struct {
	//HttpRequestsCounter    *prometheus.CounterVec
	HttpLatenciesHistogram *prometheus.HistogramVec
	Service                Service
}

func (card Resource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/card").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well

	tags := []string{"card"}

	ws.Route(ws.POST("/auth").
		To(card.middleware(card.Service.WrapMiddlewares(card.Service.Auth))).
		//To(card.auth).
		// docs
		Doc("card auth").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(Request{}).
		Writes(Response{}).
		Returns(200, "OK", Response{}).
		Returns(500, "server error", nil))

	ws.Route(ws.POST("/capture").
		To(card.middleware(card.Service.WrapMiddlewares(card.Service.Capture))).
		//To(card.capture).
		// docs
		Doc("card capture based on previous auth").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(Request{}).
		Writes(Response{}).
		Returns(200, "OK", Response{}).
		Returns(500, "server error", nil))

	ws.Route(ws.POST("/refund").
		To(card.middleware(card.Service.WrapMiddlewares(card.Service.Refund))).
		//To(card.refund).
		// docs
		Doc("card refund based on previous capture").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(Request{}).
		Writes(Response{}).
		Returns(200, "OK", Response{}).
		Returns(500, "server error", nil))
	return ws
}

func (card *Resource) middleware(handler handlerFunc) func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		var err error
		start := time.Now()
		// counter inc/ histogram observe
		defer func() {
			elapsed := float64(time.Since(start) / time.Millisecond)
			if err != nil {
				//card.HttpRequestsCounter.WithLabelValues("500", "POST", request.Request.RequestURI).Inc()
				card.HttpLatenciesHistogram.WithLabelValues("500", "POST", request.Request.RequestURI).Observe(elapsed)
			} else {
				//card.HttpRequestsCounter.WithLabelValues("200", "POST", request.Request.RequestURI).Inc()
				card.HttpLatenciesHistogram.WithLabelValues("200", "POST", request.Request.RequestURI).Observe(elapsed)
			}
		}()

		cardRequest := new(Request)
		err = request.ReadEntity(&cardRequest)
		if err != nil {
			_ = response.WriteError(http.StatusInternalServerError, err)
			return
		}
		var cardResponse *Response
		cardResponse, err = handler(cardRequest)
		if err != nil {
			_ = response.WriteError(http.StatusInternalServerError, err)
			return
		}

		_ = response.WriteEntity(cardResponse)
	}
}

// equals to...
//func (card *Resource) auth(request *restful.Request, response *restful.Response) {
//	cardRequest := new(Request)
//	err := request.ReadEntity(&cardRequest)
//	if err != nil {
//		_ = response.WriteError(http.StatusInternalServerError, err)
//		return
//	}
//	cardResponse, err := card.Service.Auth(cardRequest)
//	if err != nil {
//		_ = response.WriteError(http.StatusInternalServerError, err)
//		return
//	}
//
//	_ = response.WriteEntity(cardResponse)
//}
//
//func (card *Resource) capture(request *restful.Request, response *restful.Response) {
//
//}
//
//func (card *Resource) refund(request *restful.Request, response *restful.Response) {
//
//}
