package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, ic istioer, apiToken string) {
	authorized := r.Group("/")
	authorized.Use(authorizedWithToken(apiToken))
	authorized.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	authorized.GET("/api/routes", func(c *gin.Context) {
		var routes, err = ic.listRegisteredRoutes()
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get all routes",
			})
		}
		c.JSON(200, routes)
	})

	authorized.POST("/api/routes/*path", func(c *gin.Context) {
		path := c.Param("path")
		var r *route
		var err error
		if r, err = unmarshalRoute(path, c.Request.Body); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "improper json body",
			})
		}
		if err := ic.createVirtualService(*r); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "error processing route creation request",
			})
		}
		c.Status(http.StatusCreated)
	})

	authorized.DELETE("/api/routes/*path", func(c *gin.Context) {
		path := c.Param("path")
		if err := ic.deleteRoute(getVirtualServiceNameWithPrefix(path)); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete route",
			})
		}
		c.Status(http.StatusNoContent)
	})
}

func authorizedWithToken(apiToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.Request.Header.Get("Authorization")
		// Not bearer :o
		bearerSplit := strings.Split(reqToken, "token ")
		if len(bearerSplit) == 2 && bearerSplit[1] == apiToken {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "invalid credentials",
		})

	}
}

func validateRequired(paramName string, paramValue string) error {
	if paramValue == "" {
		return fmt.Errorf("missing required param %s", paramName)
	}
	return nil
}
