package keyvault

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/Azure/go-autorest/autorest"
)

func (k *Keyvault) logger() {
	setDebug := strings.ToLower(os.Getenv("DEBUG"))
	if strings.ToLower(setDebug) == "true" {
		k.basicClient.RequestInspector = logRequest()
		k.basicClient.ResponseInspector = logResponse()
	}
}

func logRequest() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpRequestOut(r, true)
			log.Println(string(dump))
			return r, err
		})
	}
}

func logResponse() autorest.RespondDecorator {
	return func(p autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(r *http.Response) error {
			err := p.Respond(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpResponse(r, true)
			log.Println(string(dump))
			return err
		})
	}
}
