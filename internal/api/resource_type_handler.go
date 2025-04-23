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
	var req struct {
		Name        *string `json:"name" binding:"required,min=1,max=50"`
		Description *string `json:"description" binding:"omitzero,max=255"`
	}

	if err := readJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resourceType := &store.ResourceType{
		Name:        *req.Name,
		Description: *req.Description,
	}

	resourceType, err := rh.ResourceTypeStore.CreateResourceType(resourceType)
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

	var req struct {
		Name        *string `json:"name" binding:"omitzero,max=50"`
		Description *string `json:"description" binding:"omitzero,max=255"`
	}

	if err := readJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		resourceType.Name = *req.Name
	}

	if req.Description != nil {
		resourceType.Description = *req.Description
	}

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
