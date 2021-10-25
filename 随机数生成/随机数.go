package 随机数生成

// rand5()   已知rand13(),生成rand5()
func rand5() int {
	for {
		/*
			rand13() 如果生成 1-5 return
			可选数字 (5/13) * 选中特定概率数字 (1/5) == 1/13
			如果数字范围在 8-10 重新执行 (8/13) * (1/8) == 1/13
		*/
		temp := rand13()
		if temp <= 5 {
			return temp
		}
	}
}

// rand13()   已知rand5(),生成rand13()
func rand13() int {
	for {
		/*
			rand5() 生成的随机数范围是 [1...5]
			(rand5 - 1) * 5 + rand5() 可以等概率的生成的随机数范围是 [1, 5*5]
			可以选择在 [1...13] 范围内的随机数返回。
		*/
		x := rand5()
		y := rand5()
		temp := (x-1)*5 + y
		if temp <= 13 {
			return temp
		}
		x = temp - 13 // x=12 rand12()-->rand13()
		y = rand5()
		temp = (x-1)*5 + y //60
		if temp <= 52 {    //60取 13*4 =52 个
			return 1 + (temp-1)%13
		}
	}
}
