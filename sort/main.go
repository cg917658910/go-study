package main

func insertionSort(nums []int) []int {
	for i := 1; i < len(nums); i++ {
		key := nums[i]
		j := i - 1
		for j >= 0 && nums[j] > key {
			nums[j+1] = nums[j]
			j--
		}
		nums[j+1] = key
	}
	return nums
}

func main() {
	nums := []int{5, 2, 9, 1, 5, 6}
	sorted := insertionSort(nums)
	for _, num := range sorted {
		println(num)
	}
}
