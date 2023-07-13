package contact

import (
	"bank-service/src/environments/client/resources/interfaces"
	controller "bank-service/src/libs/controllers/client"
	"bank-service/src/libs/dto"
	httpUtils "bank-service/src/libs/http"
	"bank-service/src/libs/i18n"
	"bank-service/src/utils/constant"
	"bank-service/src/utils/helpers"
	"bank-service/src/utils/pagination"

	"bank-service/src/utils/querystring"
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
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)
	pagination, err := pagination.GetPaginationFromQuery(request.URL.Query())
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	filterContacts := &dto.FilterContacts{}
	if err := querystring.Decode(filterContacts, request.URL.Query()); err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	filterContacts.UserID = jwtContext.UserID

	contacts, err := c.sContact.IndexByUserID(filterContacts, pagination)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakePaginateResponse(response, contacts, http.StatusOK, pagination)
}
func (c *contactController) Create(response http.ResponseWriter, request *http.Request) {
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)
	createContact := &dto.CreateContact{}
	err := httpUtils.GetBodyRequest(request, createContact)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	createContact.UserID = jwtContext.UserID

	contact, err := c.sContact.Create(createContact)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, contact, http.StatusCreated, i18n.T(i18n.Message{MessageID: "CONTACT.CREATED"}))

}
func (c *contactController) Update(response http.ResponseWriter, request *http.Request) {
	contactID, err := httpUtils.GetParamRequestInt(request, "id")
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)
	updateContact := &dto.UpdateContact{}
	err = httpUtils.GetBodyRequest(request, updateContact)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	updateContact.UserID = jwtContext.UserID
	updateContact.ContactID = contactID

	contact, err := c.sContact.Update(updateContact)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, contact, http.StatusOK, i18n.T(i18n.Message{MessageID: "CONTACT.UPDATED"}))
}

func (c *contactController) Delete(response http.ResponseWriter, request *http.Request) {
	contactID, err := httpUtils.GetParamRequestInt(request, "id")
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)
	filterOneContact := &dto.FilterOneContact{
		UserID:    jwtContext.UserID,
		ContactID: contactID,
	}

	contact, err := c.sContact.Delete(filterOneContact)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, contact, http.StatusOK, i18n.T(i18n.Message{MessageID: "CONTACT.DELETED"}))
}

func (c *contactController) GetOne(response http.ResponseWriter, request *http.Request) {
	contactID, err := httpUtils.GetParamRequestInt(request, "id")
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)
	filterOneContact := &dto.FilterOneContact{
		UserID:    jwtContext.UserID,
		ContactID: contactID,
	}

	contact, err := c.sContact.GetOne(filterOneContact)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, contact, http.StatusOK, i18n.T(i18n.Message{MessageID: "CONTACT.FOUND"}))
}
