/*
Copyright 2020 Splunk Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes bootstraps http routes
func RegisterRoutes(r *gin.Engine, ic Istioer, apiToken string) {
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
		if err := ic.deleteRoute(path); err != nil {
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
		// Ignore /ping handler
		if c.Request.URL.Path == "/ping" {
			c.Next()
			return
		}

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
