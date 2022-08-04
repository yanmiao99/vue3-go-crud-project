package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
)

/*

gin
gorm
mysql

*/
func main() {

	// 如何链接数据库

	dsn := "root:12345678@tcp(127.0.0.1:3306)/crud-list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 解决查表的时候会自动添加复数的问题 , 例如 user 变成了 users
			SingularTable: true,
		},
	})

	fmt.Println(db)
	fmt.Println(err)

	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 10 秒钟

	// 结构体
	type List struct {
		gorm.Model
		Name    string `gorm:"type:varchar(20); not null" json:"name" binding:"required"`
		State   string `gorm:"type:varchar(20); not null" json:"state" binding:"required"`
		Phone   string `gorm:"type:varchar(20); not null" json:"phone" binding:"required"`
		Email   string `gorm:"type:varchar(40); not null" json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200); not null" json:"address" binding:"required"`
	}

	/* 注意点 :
	1. 	结构体里面的变量 (Name) 必须是首字符大写
	gorm 指定类型
	json 表示json接受的时候的名称
	binding required 表示必须传入
	*/

	// 1. 主键没有  (不符合规范 )  给结构体添加  gorm.Model
	// 2. 名称变成复数的问题

	db.AutoMigrate(&List{})

	// 接口
	r := gin.Default()

	// 测试

	//r.GET("/", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "请求成功",
	//	})
	//})

	/* 业务码约定

	正确 : 200
	错误 : 400

	*/

	// 增

	r.POST("/user/add", func(c *gin.Context) {
		var data List

		err := c.ShouldBindJSON(&data)

		// 判断绑定是否有错误
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			// 数据库的操作
			db.Create(&data) // 创建一条数据

			c.JSON(200, gin.H{
				"msg":  "添加成功",
				"data": data,
				"code": 200,
			})
		}

	})

	// 删
	// 1. 找到对应的 id 所对应的条目
	// 2. 判断 id 是否存在
	// 3. 从数据库中删除
	// 3. 返回, id 没有找到

	// restful 编码规范 风格

	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []List

		// 接受 id
		id := c.Param("id")

		// 判断 id 是否存在
		db.Where("id = ?", id).Find(&data)

		// id 存在的情况, 则删除, 不存在则报错
		if len(data) == 0 {
			c.JSON(200, gin.H{
				"msg":  "id没有找到, 删除失败",
				"code": 400,
			})
		} else {
			// 操作数据库删除
			db.Where("id = ?", id).Delete(&data)

			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}
	})

	// 改

	r.PUT("/user/update/:id", func(c *gin.Context) {

		// 1. 找到对应的 id 所对应的条目
		// 2. 判断 id 是否存在
		// 3. 修改对应条目
		// 4. 返回 id,没有找到

		var data List

		// 接受 id
		id := c.Param("id")

		// 判断 id 是否存在
		db.Select("id").Where("id = ? ", id).Find(&data)

		// 判断 id 是否存在
		if data.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户id没有找到",
				"code": 400,
			})
		} else {
			err := c.ShouldBindJSON(&data)

			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {
				// db 修改数据库内容
				db.Where("id = ?", id).Updates(&data)

				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		}
	})

	// 查 ( 条件查询 , 全部查询 / 分页查询)

	// 条件查询
	r.GET("/user/list/:name", func(c *gin.Context) {

		// 获取路径参数
		name := c.Param("name")

		var dataList []List

		// 查询数据库
		db.Where("name = ? ", name).Find(&dataList)

		// 判断是否查询到数据
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": dataList,
			})
		}
	})

	// 全部查询
	r.GET("/user/list", func(c *gin.Context) {

		var dataList []List

		// 1. 查询全部数据,  查询分页数据
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))

		// 判断是否需要分页
		if pageSize == 0 {
			pageSize = -1
		}
		if pageNum == 0 {
			pageNum = -1
		}

		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}

		// 返回一个总数
		var total int64
		// 查询数据库
		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)

		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}

	})

	// 端口号
	PORT := "3000"
	r.Run(":" + PORT)

}
