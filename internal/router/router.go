package router

import (
	"github.com/gin-gonic/gin"
	"mastcat/internal/handler"
	"mastcat/internal/middleware"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(middleware.ErrorHandler())    //错误处理中间件
	r.Use(middleware.LimitMiddleware()) //限流中间件
	r.Use(middleware.CrawlerFilter())   //爬虫过滤中间件
	//	r.Use(middleware.AccessLogger())    //访问日志中间件
	r.POST("api/register", handler.Register)
	r.POST("api/login", handler.Login)

	authorized := r.Group("/") //权限验证中间件
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("api/blogs/manger", handler.GetBlogManger)
		authorized.POST("api/blogs", handler.AddBlog)
		authorized.PUT("api/blogs/:id", handler.UpdateBlog)
		authorized.DELETE("api/blogs/:id", handler.DeleteBlog)
		authorized.POST("api/categories", handler.AddCategory)
		authorized.PUT("api/categories/:id", handler.UpdateCategory)
	}

	r.GET("api/blogs", handler.GetBlogs)
	r.GET("api/categories", handler.GetCategories)
	r.GET("api/blogs/category/:categoryID", handler.GetBlogsByCategory)
	r.GET("api/blogs/:id", handler.GetBlogByID)
	return r
}
