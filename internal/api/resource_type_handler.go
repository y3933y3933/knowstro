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
	resourceTypeStore store.ResourceTypeStore
	logger            *slog.Logger
}

func NewResourceTypeHandler(resourceTypeStore store.ResourceTypeStore, logger *slog.Logger) *ResourceTypeHandler {
	return &ResourceTypeHandler{
		resourceTypeStore: resourceTypeStore,
		logger:            logger,
	}
}

func (rh *ResourceTypeHandler) ListTypes(c *gin.Context) {
	types, err := rh.resourceTypeStore.GetAllResourceType()
	if err != nil {
		rh.logger.Error(err.Error())
		response.InternalError(c)
		return
	}

	response.SuccessOK(c, types)
}

func (rh *ResourceTypeHandler) CreateType(c *gin.Context) {
	var req struct {
		Name        *string `json:"name" binding:"required,min=1,max=50"`
		Description *string `json:"description" binding:"omitzero,max=255"`
	}

	if err := utils.ReadJSON(c, &req); err != nil {
		rh.logger.Error(err.Error())
		if details, isValid := utils.ValidationErrors(err); !isValid {
			response.FailedValidationError(c, details)
		} else {
			response.BadRequest(c, err.Error())
		}
		return

	}

	resourceType := &store.ResourceType{}
	if req.Name != nil {
		resourceType.Name = *req.Name
	}

	if req.Description != nil {
		resourceType.Description = *req.Description
	}

	resourceType, err := rh.resourceTypeStore.CreateResourceType(resourceType)
	if err != nil {
		rh.logger.Error(err.Error())
		switch {
		case errors.Is(err, store.ErrDuplicateResourceType):
			response.UnprocessableError(c, err.Error())

		default:
			response.InternalError(c)
		}

		return
	}

	response.SuccessOK(c, resourceType)

}

func (rh *ResourceTypeHandler) UpdateType(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
	if err != nil {
		rh.logger.Error(err.Error())
		response.RecordNotFound(c)
		return
	}

	resourceType, err := rh.resourceTypeStore.GetResourceTypeByID(id)
	if err != nil {
		rh.logger.Error(err.Error())
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
		rh.logger.Error(err.Error())
		if details, isValid := utils.ValidationErrors(err); !isValid {
			response.FailedValidationError(c, details)
		} else {
			response.BadRequest(c, err.Error())
		}
		return
	}

	if req.Name != nil {
		resourceType.Name = *req.Name
	}

	if req.Description != nil {
		resourceType.Description = *req.Description
	}

	_, err = rh.resourceTypeStore.UpdateResourceType(resourceType)
	if err != nil {
		rh.logger.Error(err.Error())
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

	response.SuccessOK(c, resourceType)

}

func (rh *ResourceTypeHandler) GetTypeByID(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
	if err != nil {
		rh.logger.Error(err.Error())
		response.RecordNotFound(c)
		return
	}

	resourceType, err := rh.resourceTypeStore.GetResourceTypeByID(id)
	if err != nil {
		rh.logger.Error(err.Error())
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.RecordNotFound(c)
		default:
			response.InternalError(c)
		}
		return
	}

	response.SuccessOK(c, resourceType)

}

func (rh *ResourceTypeHandler) DeleteType(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
	if err != nil {
		rh.logger.Error(err.Error())
		response.RecordNotFound(c)
		return
	}

	err = rh.resourceTypeStore.DeleteResourceType(id)
	if err != nil {
		rh.logger.Error(err.Error())
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.RecordNotFound(c)
		default:
			response.InternalError(c)
		}
	}
	response.SuccessOK(c, nil)
}

func (rh *ResourceTypeHandler) ResetTypes(c *gin.Context) {
	err := rh.resourceTypeStore.ResetResourceType()
	if err != nil {
		rh.logger.Error(err.Error())
		response.InternalError(c)
		return
	}
	response.SuccessOK(c, nil)
}
