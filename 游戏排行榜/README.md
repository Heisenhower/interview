##### 一、题目需求
    1. 每月活动玩家得到的总分数 0-10000。
    2. 每月活动结束，根据总分是从高到底建立排行榜。
    3. 玩家分数相同，按分数排序，先到先得。
    4. 玩家能查询自己名次前后十位玩家的分数和名次。

***
***
##### 二、解决思路
    1、对于10万用户量的程序来说，使用MySQL就能解决此问题，但是对于排榜使用Redis才是最优解，
    单台Redis性能越能达到7.5万/秒的处理能力，在时间和效率来说是最优解。
    Redis 中的有序集合(sorted set)能够处理这个需求，利用其成员的唯一性、有序性。


在redis 中使用zset 建立 user_score   score + id

~~~
/*
@Parma setKey user_info
@Parma id 玩家id
@Parma score 玩家游戏分数
*/

// addUserInfo 向redis 中添加游戏玩家成绩
func addUserInfo(setKey, id string, score int64) {
	s := scoreDeviation(score)
	setInfo := []redis.Z{
		redis.Z{
			Score:  float64(s),
			Member: id,
		},
	}
	_, err := rdb.ZAdd(setKey, setInfo...).Result()
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
}
~~~
***

2、通过 ZREVRANGE 查询出全部玩家的成绩排行

`ZREVRANGE user_score 0 -1 WITHSCORES`

~~~
/*
@Parma key user_info
@Parma m 全部玩家排名和成绩
*/
// getRankAll 获取排行榜总榜（降序）
func getRankAll(key string) (m map[interface{}]interface{}) {
	ret, err := rdb.ZRevRangeWithScores(key, 0, -1).Result()
	if err != nil {
		fmt.Printf("zrevrange failed, err:%v\n", err)
		return
	}
	m = make(map[interface{}]interface{})
	for _, z := range ret {
		m[z.Member] = scoreDeduct(int64(z.Score))
	}
	return m
}
~~~
***
3、为满足同分数先到先得，在原始游戏分数上加上时间和偏移量，来确保同分排名有先后


~~~
// scoreDeviation 游戏得分加偏移量
func scoreDeviation(score int64) (count int64) {
	count = int64(math.Pow(10, 14))*score + 9999999999999 - time.Now().UnixNano()/1e6
	return count
}

//玩家分数 * 偏移量(int64(math.Pow(10, 14))) + 最大时间（9999999999999 ）13位 -  时间戳（time.Now().UnixNano()/1e6）

// scoreDeviation 游戏得分去除偏移量
func scoreDeduct(score int64) (count int64) {
	count = score / 1e14
	return count
}
~~~

***
4、玩家查询自己的排名，和自己前后10名的成绩

    先通过id查询当前玩家排名 x

`ZREVRANK user_score 当前玩家id`

    再根据当前玩家排名 x 去查询前后10名玩家id

`ZREVRANGE user_score x-10 x+10`

~~~
/*
@Parma key user_info
@Parma id 当前玩家id
@Parma interval 成绩间隔 10
@Parma rank 当前玩家排名 
@Parma subScoreRankMap 前后10名玩家分数和成绩 
@
*/
// getUserRank 玩家查询名次 包含范围
func getUserRank(key, id string, interval int64) (rank int64, subScoreRankMap map[string]string) {
	ret, err := rdb.ZRevRank(key, id).Result()
	if err != nil {
		fmt.Printf("ZScore failed, err:%v\n", err)
		return
	}
	//根据当前玩家名次取查询其前后interval位玩家的分数和名次
	arr, err := rdb.ZRevRange(key, ret-interval, ret+interval).Result()
	if err != nil {
		fmt.Printf("ZRevRange failed, err:%v\n", err)
		return
	}
	subScoreRankMap = make(map[string]string)
	for _, id := range arr {
		rank, err := rdb.ZRevRank(key, id).Result()
		if err != nil {
			fmt.Printf("ZScore failed, err:%v\n", err)
			continue
		}
		score, err := rdb.ZScore(key, id).Result()
		if err != nil {
			fmt.Printf("ZScore failed, err:%v\n", err)
			continue
		}
		subScoreRankMap[id] = fmt.Sprint(rank, ":", scoreDeduct(int64(score)))
	}
	return ret, subScoreRankMap
}
~~~