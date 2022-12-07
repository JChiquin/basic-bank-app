package notfound

import (
	"net/http"

	myErrors "bank-service/src/libs/errors"
	utilsHttp "bank-service/src/libs/http"

	"github.com/gorilla/mux"
)

//CustomNotFoundHanlder overrides the default not found handler on router
func CustomNotFoundHandler(muxRouter *mux.Router) {
	muxRouter.NotFoundHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		utilsHttp.MakeErrorResponse(response, myErrors.ErrURLNotFound)
	})
}
