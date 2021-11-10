package v1

//import (
//	"backend/internal/service/db"
//	"github.com/labstack/echo/v4"
//	"net/http"
//)
//
//func (h *Handler) initLinksRoutes(api *echo.Group) {
//	links := api.Group("/links")
//	{
//		links.GET("/", h.GetListLink)
//		links.GET("/:id", h.GetLink)
//		links.POST("/", h.AddLink)
//		links.PUT("/:id", h.UpdateLink)
//	}
//}
//
//type Links struct {
//	Urls *[]db.ResponseLink `json:"Links"`
//}
//
//func (h Handler) GetListLink(c echo.Context) error {
//	links, err := h.Services.DB.GetAllLinks()
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	}
//
//	return c.JSON(http.StatusOK, Links{Urls: links})
//}
//
//func (h Handler) GetLink(c echo.Context) error {
//	id := c.Param("id")
//
//	link, err := h.Services.DB.GetLinkById(id)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	}
//
//	return c.JSON(http.StatusOK, link)
//}
//
//func (h Handler) AddLink(c echo.Context) error {
//	var link db.Link
//	if err := c.Bind(&link); err != nil {
//		return c.JSON(http.StatusBadRequest, err)
//	}
//
//	id, err := h.Services.DB.AddLink(link)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	}
//
//	err = h.Services.Queue.SetLinkStatus(id, link.URL)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	}
//
//	return c.JSON(http.StatusOK, struct {
//		Id int `json:"id"`
//	}{Id: id})
//}
//
//func (h Handler) UpdateLink(c echo.Context) error {
//	req := &struct {
//		Status int `json:"status_code"`
//	}{}
//
//	id := c.Param("id")
//
//	if err := c.Bind(&req); err != nil {
//		return c.JSON(http.StatusBadRequest, err)
//	}
//
//	err := h.Services.DB.UpdateStatusLink(req.Status, id)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	}
//
//	return c.NoContent(http.StatusNoContent)
//}
