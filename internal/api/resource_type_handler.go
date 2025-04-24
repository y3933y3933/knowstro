package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/store"
	"github.com/y3933y3933/knowstro/internal/utils"
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
	types, err := rh.ResourceTypeStore.GetAllResourceType()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceTypes": types,
	})
}

func (rh *ResourceTypeHandler) CreateType(c *gin.Context) {
	var req struct {
		Name        *string `json:"name" binding:"required,min=1,max=50"`
		Description *string `json:"description" binding:"omitzero,max=255"`
	}

	if err := utils.ReadJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resourceType := &store.ResourceType{}
	if req.Name != nil {
		resourceType.Name = *req.Name
	}

	if req.Description != nil {
		resourceType.Description = *req.Description
	}

	resourceType, err := rh.ResourceTypeStore.CreateResourceType(resourceType)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateResourceType):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong."})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceType": resourceType,
	})

}

func (rh *ResourceTypeHandler) UpdateType(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
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

	if err := utils.ReadJSON(c, &req); err != nil {
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
		case errors.Is(err, store.ErrDuplicateResourceType):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"resourceType": resourceType})

}

func (rh *ResourceTypeHandler) GetTypeByID(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
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

	c.JSON(http.StatusOK, gin.H{"resourceType": resourceType})

}

func (rh *ResourceTypeHandler) DeleteType(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = rh.ResourceTypeStore.DeleteResourceType(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
	}
}

func (rh *ResourceTypeHandler) ResetTypes(c *gin.Context) {
	err := rh.ResourceTypeStore.ResetResourceType()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset resource types"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "resource types reset successfully"})
}
