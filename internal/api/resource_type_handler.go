package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/store"
)

type ResourceTypeHandler struct {
	ResourceTypeStore store.ResourceTypeStore
}

type ResourceTypeRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}

func NewResourceTypeHandler(resourceTypeStore store.ResourceTypeStore) *ResourceTypeHandler {
	return &ResourceTypeHandler{
		ResourceTypeStore: resourceTypeStore,
	}
}

func (rh *ResourceTypeHandler) ListTypes(c *gin.Context) {

	c.JSON(http.StatusOK, map[string]any{
		"id":          1,
		"name":        "後端",
		"description": "後端",
	})
}

func (rh *ResourceTypeHandler) CreateType(c *gin.Context) {
	var req ResourceTypeRequest

	if err := readJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resourceType, err := rh.ResourceTypeStore.CreateResourceType(&store.ResourceType{Name: req.Name})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resourceType,
	})

}

func (rh *ResourceTypeHandler) UpdateType(c *gin.Context) {
	id, err := readIDParam(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	resourceType, err := rh.ResourceTypeStore.GetResourceTypeByID(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	var req ResourceTypeRequest

	if err := readJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resourceType.Name = req.Name
	_, err = rh.ResourceTypeStore.UpdateResourceType(resourceType)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resourceType})

}

func (rh *ResourceTypeHandler) GetTypeByID(c *gin.Context) {
	id, err := readIDParam(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	resourceType, err := rh.ResourceTypeStore.GetResourceTypeByID(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resourceType})

}
