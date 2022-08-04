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
业务码规范 :
请求成功 : 200
请求失败 : 400
*/

//gin文档  : https://www.kancloud.cn/shuangdeyu/gin_book/949411
//gorm文档 : https://gorm.io/zh_CN/docs/

func main() {
	// 链接到数据库
	dsn := "root:12345678@tcp(127.0.0.1:3306)/go-crud-list?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 解决查表的时候会自动添加复数的问题 , 例如 user 变成了 users
			SingularTable: true,
		},
	})
	//fmt.Println(err)
	//fmt.Println(db)
	// 配置数据库
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, _ := db.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 10秒钟

	// 结构体
	type List struct {
		gorm.Model
		Name    string `gorm:"type: varchar(20);not null" json:"name" binding:"required"`
		State   string `gorm:"type: varchar(20);not null" json:"state" binding:"required"`
		Phone   string `gorm:"type: varchar(20);not null" json:"phone" binding:"required"`
		Email   string `gorm:"type: varchar(40);not null" json:"email" binding:"required"`
		Address string `gorm:"type: varchar(200);not null" json:"address" binding:"required"`
	}

	// State 0 在职 1 试用期

	// 注意点 :
	// 结构体里面的变量必须大写 , 首字母大写是表示可公开访问
	// gorm: 指定类型 , 并且不能为空
	// json: 指定json接受的时候的名称
	// binding:"required" 如果加上了这个 , 但是传入值没有必填, 则会直接报错

	// 迁移表
	_ = db.AutoMigrate(&List{})

	// 接口
	r := gin.Default()

	/**
	 * @name   /user/list
	 * @param  pageSize  返回多少条数据
	 * @param  pageNum   第几页
	 * @return 所有的列表数据, 或者分页数据
	 */
	r.GET("/user/list", func(c *gin.Context) {

		var listRes []List // 需要使用切片类型获取数据, 否则只能获取到一条数据

		// 查数据库, 查全部数据
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))

		// 判断是否需要分页
		if pageSize == 0 {
			pageSize = -1 // -1 代表不使用 limit
		}

		if pageNum == 0 {
			pageNum = -1 // -1 代表不使用 offset
		}

		// 处理分页
		offset := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offset = -1
		}

		// 获取总数
		var total int64

		// 查询数据库
		// Model(结构体).Count(总数).Offset(分页页数).Limit(页数).Find(查询且数据插入到哪)
		db.Model(listRes).Count(&total).Limit(pageSize).Offset(offset).Find(&listRes)

		// 判断是否有返回的数据
		if len(listRes) == 0 {
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
					"list":     listRes,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}
	})

	/**
	 * @name    /user/list:name
	 * @param   name  名称
	 * @return  查询单个列表 (查询到多个, 也会返回多个)
	 */

	r.GET("/user/list/:name", func(c *gin.Context) {
		name := c.Param("name") // 获取路径中的参数
		//name := c.Query("name") // 获取 name="name" 传递的参数

		var listRes []List // 必须使用切片, 否则只会返回一个

		// 查询数据库
		db.Where("name = ?", name).Find(&listRes)

		// 判断是否有返回的数据
		if len(listRes) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": listRes,
			})
		}
	})

	/**
	 * @name    /user/add
	 * @param   {
	    "name":"测试数据",
	    "state": "0",
	    "phone": "13800138000",
	    "email": "13800138000@qq.com",
	    "address": "广东省中山市东区"
		}
	 * @return  添加一条列表
	*/
	r.POST("/user/add", func(c *gin.Context) {

		var listRes List // 因为是添加一条, 所以直接定义即可

		// ShouldBindJSON 表示绑定 json
		err := c.ShouldBindJSON(&listRes)

		// 没有错误就会返回 nil
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			db.Create(&listRes) // 通过数据的指针来创建
			c.JSON(200, gin.H{
				"msg":  "添加成功",
				"code": 200,
				"data": listRes,
			})
		}
	})

	/**
	 * @name    /user/update/:id
	 * @param   id  条目的id
	 * @return  修改一条列表信息
	 */

	r.PUT("/user/update/:id", func(c *gin.Context) {

		var listRes List // 因为是添加一条, 所以直接定义即可

		// 获取需要编辑的条目的id
		id := c.Param("id")

		// 查找对应的 id 是否存在
		db.Select("id").Where("id = ?", id).Find(&listRes)

		// 先判断 id 是否存在
		if listRes.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户ID不存在",
				"code": 400,
			})
		} else {
			// ShouldBindJSON 表示绑定 json
			err := c.ShouldBindJSON(&listRes)

			// 修改数据
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {
				db.Where("id = ?", id).Updates(&listRes) // 更新数据
				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		}
	})

	/**
	 * @name    /user/delete/:id
	 * @param   id  条目的id
	 * @return  删除一条列表信息
	 */
	r.DELETE("/user/delete/:id", func(c *gin.Context) {

		var listRes []List

		// 获取 id
		id := c.Param("id")

		// 查找对应的 id 是否存在
		db.Where("id = ?", id).Find(&listRes)

		// id 存在,则进行删除
		if len(listRes) == 0 {
			c.JSON(200, gin.H{
				"msg":  "删除失败",
				"code": 400,
			})
		} else {
			// 删除数据库的 id 对应值 (软删除 , 在数据库总不会真正被删除)
			db.Where("id = ? ", id).Delete(&listRes)
			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}
	})

	// 启动端口号为 3000
	PORT := "3000"
	fmt.Println("服务器启动成功 - 端口号 : " + PORT)
	r.Run(":" + PORT)
}
