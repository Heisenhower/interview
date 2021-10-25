package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const PORT = ":8080" //web 服务端口

// GetRandom 获取浮点随机数
func GetRandom(floor, ceiling float64) float64 {
	return floor + rand.Float64()*(ceiling-floor)

}

// ShuffleCardsLogic 进行洗牌操作
func ShuffleCardsLogic(pokerSlice []int) (arr []int) {
	length := len(pokerSlice)
	rand.Seed(time.Now().UnixNano())
	for i := length - 1; i > 0; i-- {
		j := int(GetRandom(0, float64(i)+1))
		//每次随机选出一个元素置于最后，然后除去最后这个元素之外的序列依次递归
		//只要保证每次“置后”是真正的随机选择，那么最终的序列就是原始序列的一个随机排列。
		pokerSlice[i], pokerSlice[j] = pokerSlice[j], pokerSlice[i]
	}
	fmt.Println(pokerSlice)
	return pokerSlice
}

// ShuffleCardsHandler 处理洗牌请求
func ShuffleCardsHandler(c *gin.Context) {
	var pokerSliceInt []int
	pokers := c.PostForm("original")
	if pokers == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "original pokers is null",
		})
		return
	}
	pokerSlice := strings.Split(pokers, ",")
	for _, item := range pokerSlice {
		val, err := strconv.Atoi(item)
		if err != nil {
			fmt.Println("slice change failed")
			continue
		}
		pokerSliceInt = append(pokerSliceInt, val)
	}
	//输出json结果给调用方
	c.JSON(http.StatusOK, gin.H{
		"shuffleCards": ShuffleCardsLogic(pokerSliceInt),
		"original":     pokers,
	})

}

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/shuffleCards", ShuffleCardsHandler)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	return r
}

func main() {
	r := SetupRouter()
	if err := r.Run(PORT); err != nil {
		fmt.Printf("StartUp Service Failed, Err:%v\n", err)
	}
}
