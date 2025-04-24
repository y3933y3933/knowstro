package api

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/response"
	"github.com/y3933y3933/knowstro/internal/store"
	"github.com/y3933y3933/knowstro/internal/utils"
)

type ResourceTypeHandler struct {
	ResourceTypeStore store.ResourceTypeStore
	Logger            *slog.Logger
}

func NewResourceTypeHandler(resourceTypeStore store.ResourceTypeStore, logger *slog.Logger) *ResourceTypeHandler {
	return &ResourceTypeHandler{
		ResourceTypeStore: resourceTypeStore,
	}
}

func (rh *ResourceTypeHandler) ListTypes(c *gin.Context) {
	types, err := rh.ResourceTypeStore.GetAllResourceType()
	if err != nil {
		rh.Logger.Error(err.Error())
		response.InternalError(c)
		return
	}

	response.Success(c, types)

}

func (rh *ResourceTypeHandler) CreateType(c *gin.Context) {
	var req struct {
		Name        *string `json:"name" binding:"required,min=1,max=50"`
		Description *string `json:"description" binding:"omitzero,max=255"`
	}

	if err := utils.ReadJSON(c, &req); err != nil {
		rh.Logger.Error(err.Error())
		response.BadRequest(c, err.Error())
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
		rh.Logger.Error(err.Error())
		switch {
		case errors.Is(err, store.ErrDuplicateResourceType):
			response.UnprocessableError(c, err.Error())

		default:
			response.InternalError(c)
		}

		return
	}

	response.Success(c, resourceType)

}

func (rh *ResourceTypeHandler) UpdateType(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
	if err != nil {
		rh.Logger.Error(err.Error())
		response.RecordNotFound(c)
		return
	}

	resourceType, err := rh.ResourceTypeStore.GetResourceTypeByID(id)
	if err != nil {
		rh.Logger.Error(err.Error())
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.RecordNotFound(c)
		default:
			response.InternalError(c)
		}
		return
	}

	var req struct {
		Name        *string `json:"name" binding:"omitzero,max=50"`
		Description *string `json:"description" binding:"omitzero,max=255"`
	}

	if err := utils.ReadJSON(c, &req); err != nil {
		rh.Logger.Error(err.Error())
		response.BadRequest(c, err.Error())
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
		rh.Logger.Error(err.Error())
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.RecordNotFound(c)
		case errors.Is(err, store.ErrDuplicateResourceType):
			response.UnprocessableError(c, err.Error())
		default:
			response.InternalError(c)
		}
		return
	}

	response.Success(c, resourceType)

}

func (rh *ResourceTypeHandler) GetTypeByID(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
	if err != nil {
		rh.Logger.Error(err.Error())
		response.RecordNotFound(c)
		return
	}

	resourceType, err := rh.ResourceTypeStore.GetResourceTypeByID(id)
	if err != nil {
		rh.Logger.Error(err.Error())
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.RecordNotFound(c)
		default:
			response.InternalError(c)
		}
		return
	}

	response.Success(c, resourceType)

}

func (rh *ResourceTypeHandler) DeleteType(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
	if err != nil {
		rh.Logger.Error(err.Error())
		response.RecordNotFound(c)
		return
	}

	err = rh.ResourceTypeStore.DeleteResourceType(id)
	if err != nil {
		rh.Logger.Error(err.Error())
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.RecordNotFound(c)
		default:
			response.InternalError(c)
		}
	}
	response.Success(c, nil)
}

func (rh *ResourceTypeHandler) ResetTypes(c *gin.Context) {
	err := rh.ResourceTypeStore.ResetResourceType()
	if err != nil {
		rh.Logger.Error(err.Error())
		response.InternalError(c)
		return
	}
	response.Success(c, nil)
}
