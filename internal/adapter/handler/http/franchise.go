package http

import (
	"net/http"

	"github.com/OzkrOssa/franchi-api/internal/core/domain"
	"github.com/OzkrOssa/franchi-api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type FranchiseHandler struct {
	svc port.FranchiseService
}

func NewFranchiseHandler(svc port.FranchiseService) *FranchiseHandler {
	return &FranchiseHandler{svc}
}

type newFranchiseRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *FranchiseHandler) NewFranchise(ctx *gin.Context) {
	var req newFranchiseRequest

	// Validar el JSON entrante
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid JSON body",
		})
		return
	}

	franchise := domain.Franchise{
		Name: req.Name,
	}

	_, err := h.svc.CreateNewFranchise(ctx, &franchise)
	if err != nil {
		switch err {
		case domain.ErrConflictingData:
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "conflict",
				"message": "Franchise already exists",
			})
			return
		case domain.ErrInvalidData:
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "invalid_data",
				"message": "name must be valid",
			})
			return
		case domain.ErrInternal:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Could not create franchise",
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Franchise created successfully",
	})
}

func (h *FranchiseHandler) UpdateFranchise(ctx *gin.Context) {
	var req newFranchiseRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	franchise := domain.Franchise{
		Name: req.Name,
	}

	_, err := h.svc.CreateNewFranchise(ctx, &franchise)

	if err != nil {
		switch err {
		case domain.ErrConflictingData:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err,
			})
		case domain.ErrInternal:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "franchise created successfully",
	})
}
