package main

import "fmt"

/**
 * 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
 *
 *
 * @param num int整型一维数组
 * @param size int整型
 * @return int整型一维数组
 */
func maxInWindows(num []int, size int) []int {
	// write code here
	n := len(num)
	if size > n || size <= 0 {
		return nil
	}

	q := []int{}
	ans := []int{}
	for i := 0; i < n; i++ {
		for len(q) > 0 && num[q[len(q)-1]] < num[i] {
			q = q[:len(q)-1]
		}
		q = append(q, i)
		if i >= size && q[0] <= i-size {
			q = q[1:]
		}
		if i >= size-1 {
			ans = append(ans, num[q[0]])
		}
	}
	return ans
}

func maxSlidingWindow(nums []int, k int) []int {
	if len(nums) == 0 || k == 0 {
		return []int{}
	}

	result := []int{}
	deque := []int{} // 存的是索引，不是值！

	for i := 0; i < len(nums); i++ {
		// ❌ 1. 清除窗口外的元素（超出左边界）
		if len(deque) > 0 && deque[0] < i-k+1 {
			deque = deque[1:]
		}

		// ❌ 2. 清除比当前元素小的（他们没机会成为最大值）
		for len(deque) > 0 && nums[deque[len(deque)-1]] < nums[i] {
			deque = deque[:len(deque)-1]
		}

		// ✅ 3. 加入当前元素索引
		deque = append(deque, i)

		// 📌 4. 从第一个满窗口开始，记录最大值
		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}

	return result
}

func maxWindow(nums []int, k int) []int {
	deque := []int{} // 双端队列 索引
	ans := []int{}   // 存放结果
	for i := 0; i < len(nums); i++ {
		// 清除窗口外的元素（超出左边界）
		if len(deque) > 0 && deque[0] < i-k+1 {
			deque = deque[1:]
		}
		// 清除比当前元素小的（他们没机会成为最大值）
		for len(deque) > 0 && nums[deque[len(deque)-1]] < nums[i] {
			deque = deque[:len(deque)-1]
		}
		// 加入当前元素索引
		deque = append(deque, i)
		// 从第一个满窗口开始，记录最大值
		if i >= k-1 {
			ans = append(ans, nums[deque[0]])
		}
	}
	return ans
}

func main() {
	nums := []int{10, 3, -1, -3, 5, 3, 6, 7}
	k := 3

	ans := maxWindow(nums, k)

	fmt.Println("ans: ", ans)

}
