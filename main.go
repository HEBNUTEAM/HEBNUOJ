package main

import (
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/router"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := common.InitDB()
	defer db.Close()
	r := router.CollectRoute(gin.Default())
	panic(r.Run())
}
