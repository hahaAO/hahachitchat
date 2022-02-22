//热度统计
//计算总热度  postid-posthot
package dataLayer

import (
	"errors"
)

//总热度
type Post_idandhot struct {
	Post_id  int   `json:"post_id"`
	Post_hot int64 `json:"post_hot"`
}

const (
	timehotweight    = 1
	commenthotweight = 3600
)

//时间热度
type post_idandtimehot struct {
	post_id      int
	post_timehot int64
}

//评论热度
type post_idandcommenthot struct {
	post_id         int
	post_commenthot int64
}

//获取所有post的时间热度  post_idandtimehot(post_id和post_timehot)
func allpost_timehot() ([]post_idandtimehot, error) {
	var results []post_idandtimehot
	var err error
	aint, aposts := AllSelectPost()
	if aint == 1 { //查询成功
		along := len(aposts)
		for i := 0; i < along; i++ {
			var result post_idandtimehot
			result.post_id = aposts[i].Post_id
			result.post_timehot = timehotweight * aposts[i].Post_time.Unix()
			results = append(results, result)
		}
		return results, err
	} else if aint == 0 { //没有帖子

		return results, err
	} else { //有其他问题
		err = errors.New("allpost_timehot 有其他问题")

		return results, err
	}
}

//根据id数组 获取所有post的评论热度   post_idandcommenthot(post_id和post_commenthot)
func allpost_commenthot(post_ids []int) ([]post_idandcommenthot, error) {
	plong := len(post_ids)
	results := make([]post_idandcommenthot, plong)
	var err error
	for i := 0; i < plong; i++ {
		results[i].post_id = post_ids[i]
		aint, apcomids := AllCommentidOnpostid(results[i].post_id)
		if aint == 1 { //查询成功
			along := int64(len(apcomids))                         //评论数
			results[i].post_commenthot = along * commenthotweight //评论热度
		} else if aint == 0 { //没有评论
			results[i].post_commenthot = 0
		} else { //有其他问题
			err = errors.New("allpost_commenthot: 有其他问题")
		}
	}
	return results, err
}

//计算所有帖子对应的总热度并返回 Post_idandhot (所有帖子的post_id和总热度 post_hot),(error)
func Allposthot() ([]Post_idandhot, error) {
	var err error
	var post_ids []int //记录id的切片
	timehots, err := allpost_timehot()
	if err != nil {
		return nil, err
	}
	postlong := len(timehots)                  //帖子数量
	results := make([]Post_idandhot, postlong) //如果不初始化(分配长度)则会报错

	for i := 0; i < postlong; i++ { //初始化热度(总热度等于时间热度)
		results[i].Post_id = timehots[i].post_id        //初始化id
		results[i].Post_hot = timehots[i].post_timehot  //总热度等于时间热度
		post_ids = append(post_ids, results[i].Post_id) //记录id的切片
	}
	for i := 0; i < postlong; i++ { //总热度加上评论热度
		commenthots, err := allpost_commenthot(post_ids)
		if err != nil {

			return nil, err
		}
		results[i].Post_hot += commenthots[i].post_commenthot //总热度加上评论热度
	}
	return results, err
}
