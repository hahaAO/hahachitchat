//放一些零碎的小工具
package main

import (
	"math/rand"
	"strconv"
	"time"
)

//根据时间戳加上随机数生成唯一图片id
func timeRandId() string {
	kaishi := time.Now().UnixNano()
	timeid := strconv.FormatInt(kaishi, 10)
	randid := strconv.FormatInt(rand.Int63(), 10)
	return timeid + randid
}
