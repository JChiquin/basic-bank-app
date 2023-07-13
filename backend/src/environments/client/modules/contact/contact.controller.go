package contact

import (
	"bank-service/src/environments/client/resources/interfaces"
	controller "bank-service/src/libs/controllers/client"
	"net/http"
)

/*
struct that implements IContactController
*/
type contactController struct {
	controller.ClientController
	sContact interfaces.IContactService
}

/*
NewContactController creates a new controller, receives service by dependency injection
and returns IContactController, so it needs to implement all its methods
*/
func NewContactController(sContact interfaces.IContactService) interfaces.IContactController {
	cContact := contactController{sContact: sContact}
	return &cContact
}

func (c *contactController) Index(response http.ResponseWriter, request *http.Request) {

}
func (c *contactController) Create(response http.ResponseWriter, request *http.Request) {

}
func (c *contactController) Update(response http.ResponseWriter, request *http.Request) {

}
func (c *contactController) Delete(response http.ResponseWriter, request *http.Request) {

}
