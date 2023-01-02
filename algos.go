
func quickSort(arr []int, low, high int) []int {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func quickSortStart(arr []int) []int {
	return quickSort(arr, 0, len(arr)-1)
}

func partition(arr []int, low, high int) ([]int, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func substrgenerator() {
	arr := []string{"a", "b", "c"}
	// 	a
	// 	b
	// 	c
	// 	ab
	// 	bc
	// 	abc
	n := len(arr)
	for l := 1; l <= n; l++ {
		for i := 0; i <= n-l; i++ {
			j := i + l - 1
			out := strings.Builder{}
			for k := i; k <= j; k++ {
				out.WriteString(arr[k])
			}
			fmt.Println(out.String())
		}
	}
}
