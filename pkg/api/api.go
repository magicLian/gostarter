package api

import (
	"github.com/magicLian/gostarter/pkg/docs"
	"github.com/magicLian/gostarter/pkg/setting"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/*
当运行完 swag fmt之后，可能由于swag的bug,关于header的注释无法识别，可以复制这部分的注释替换fmt之后的header注释

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
*/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (hs *HTTPServer) apiRegister() {
	// requireSignedIn := middleware.RequireSignedIn
	// requireAdmin := middleware.RequireAdmin

	r := hs.ginEngine
	apiBasePath := "/"
	initSwaggerDocs(apiBasePath)
	r.GET("/", hs.test)
	v1Group := r.Group(apiBasePath)
	{
		// systemRouter := v1Group.Group("/system", requireAdmin)
		// {
		// 	systemRouter.GET("/info", hs.getSystemRemindRules)
		// }
		healthRouter := v1Group.Group("/health")
		{
			healthRouter.GET("/", hs.healthHandler)
		}
		testRouter := v1Group.Group("/test")
		{
			testRouter.GET("/", hs.test)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func initSwaggerDocs(apiBasePath string) {
	docs.SwaggerInfo.Title = "go starter API"
	docs.SwaggerInfo.Description = "This is api server of go starter server."
	docs.SwaggerInfo.Version = setting.APIVersion
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.BasePath = apiBasePath
}
