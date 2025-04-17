package http

import "github.com/gin-gonic/gin"

type Router struct {
	*gin.Engine
}

func NewRouter(fh *FranchiseHandler) (*Router, error) {

	router := gin.New()

	v1 := router.Group("/v1")
	{
		f := v1.Group("/franchise")
		{
			f.POST("/", fh.NewFranchise)
			f.PATCH("/", fh.UpdateFranchise)
		}
	}

	return &Router{router}, nil
}
