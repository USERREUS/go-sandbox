package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	storePath = "../storage/tasks2.json"
)

type Data struct {
	ID   int   `json:"id"`
	Prev []int `json:"prev"`
	Next []int `json:"next"`
	Time int   `json:"time"`
	Res  int   `json:"res"`
}

type Task struct {
	Prev []int
	Next []int
	Time int
	Res  int
}

func jsonParse(name string) ([]Data, error) {
	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var data []Data
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func shfl(nums []int) {
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
}

func genSample(ids []int, taskMap map[int]Task) (q []int) {
	q = append(q, ids[0])
	q = append(q, ids[1])
	used := make(map[int]bool)
	for i := 0; i < len(ids)-1; i++ {
		fmt.Println(i)
		item := q[i]
		task := taskMap[item]
		shfl(task.Next)
		for _, item := range task.Next {
			if !used[item] {
				flag := false
				for _, pr := range taskMap[item].Prev {
					if !used[pr] {
						flag = true
					}
				}
				if flag {
					continue
				}
				q = append(q, item)
				used[item] = true
			}
		}
	}
	return
}

func genSampleMany(n int, ids []int, taskMap map[int]Task) (res [][]int) {
	for i := 0; i < n; i++ {
		res = append(res, genSample(ids, taskMap))
	}
	return
}

func main() {
	//Parse JSON
	data, err := jsonParse(storePath)
	if err != nil {
		log.Fatalf("Не удалось прочитать файл: %v", err)
	}

	// Print data
	fmt.Println("Data:")
	for _, rec := range data {
		fmt.Println(rec)
	}
	fmt.Println()

	// Собираем слайс ID
	var idSlice []int
	for _, rec := range data {
		idSlice = append(idSlice, rec.ID)
	}

	// Выводим результат
	fmt.Println("idSlice:")
	fmt.Println(idSlice)
	fmt.Println()

	// Create Map of tasks
	taskMap := make(map[int]Task)
	for _, rec := range data {
		taskMap[rec.ID] = Task{
			Prev: rec.Prev,
			Next: rec.Next,
			Time: rec.Time,
			Res:  rec.Res,
		}
	}

	// Print data
	fmt.Println("taskMap:")
	for id, task := range taskMap {
		fmt.Print(id)
		fmt.Print(" : ")
		fmt.Println(task)
	}
	fmt.Println()

	//sample := genSampleMany(10, idSlice, taskMap)
	// fmt.Println("sample:")
	// for _, slice := range sample {
	// 	fmt.Println(slice)
	// }
	// fmt.Println()
	sample := genSample(idSlice, taskMap)
	fmt.Println(sample)
	fmt.Println()
}

/*
import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type Data struct {
	ID   int   `json:"id"`
	Prev []int `json:"prev"`
	Next []int `json:"next"`
	Time int   `json:"time"`
	Req  int   `json:"req"`
}

type Task struct {
	prev []int
	next []int
	time int
	req  int
}

type Pair struct {
	time int
	comp []int
}

const (
	storePath = "..storage/tasks.json"
)

func parse_json(name string) ([]Data, error) {
	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var data []Data
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func is_valid(p []int, task_map map[int]Task) bool {
	used := make(map[int]bool)
	valid := true
	for _, i := range p {
		for _, j := range task_map[i].next {
			if used[j] {
				valid = false
			}
		}
		used[i] = true
	}
	return valid
}

func calc(task_map map[int]Task, p []int, maxPers int, ch chan Pair) {
	q := make([]int, len(p))
	copy(q, p)

	done := make(map[int]bool)
	inProcess := make(map[int]int)
	time := 0
	pers := maxPers

	for len(q) > 0 {
		taskID := q[0]
		task := task_map[taskID] // replace with your actual task retrieval function

		taskReady := true
		if len(task.prev) > 0 {
			for _, i := range task.prev {
				if !done[i] {
					taskReady = false
					break
				}
			}
		}

		if taskReady && pers >= task.req {
			inProcess[taskID] = task.time
			pers -= task.req
			q = q[1:]
		} else {
			minTime := 10000
			for _, t := range inProcess {
				if t < minTime {
					minTime = t
				}
			}
			time += minTime

			for item, t := range inProcess {
				if t == minTime {
					pers += task_map[item].req
					done[item] = true
					inProcess[item] = 10000
				} else {
					inProcess[item] -= minTime
				}
			}
		}
	}

	minTime := 10000
	for _, t := range inProcess {
		if t < minTime {
			minTime = t
		}
	}
	time += minTime

	fmt.Println(time)

	ch <- Pair{time: time, comp: p}
}

func permute(nums []int, task_map map[int]Task) [][]int {
	var result [][]int
	var backtrack func(int)

	backtrack = func(first int) {
		if first >= len(nums) {
			temp := make([]int, len(nums))
			copy(temp, nums)
			if is_valid(temp, task_map) {
				result = append(result, temp)
			}
			return
		}

		for i := first; i < len(nums); i++ {
			nums[first], nums[i] = nums[i], nums[first]
			backtrack(first + 1)
			nums[first], nums[i] = nums[i], nums[first]
		}
	}

	backtrack(0)

	return result
}

func shfl(nums []int) {
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
}

// func gen_all(n int, ids []int, task_map map[int]Task) [][]int {
// 	res := make([][]int, 0)
// 	for i := 0; i < n; i++ {
// 		res = append(res, generate(ids, task_map))
// 	}
// 	return res
// }

func generate(ids []int, task_map map[int]Task) (q []int) {
	q = append(q, ids[0])

	used := make(map[int]bool)

	for i := 0; i < len(ids)-1; i++ {
		item := q[i]
		task := task_map[item]
		shfl(task.next)
		for _, item := range task.next {
			if !used[item] {
				q = append(q, item)
				used[item] = true
			}
		}
	}

	return
}

func main() {
	data, err := parse_json(storePath)
	if err != nil {
		log.Fatalf("Не удалось прочитать файл: %v", err)
	}

	// Print data
	fmt.Println(data)

	// Собираем слайс ID
	var idSlice []int
	for _, d := range data {
		idSlice = append(idSlice, d.ID)
	}

	// Выводим результат
	fmt.Println(idSlice)

	task_map := make(map[int]Task)
	for _, val := range data {
		task_map[val.ID] = Task{prev: val.Prev, next: val.Next, time: val.Time, req: val.Req}
	}

	fmt.Println(task_map)

	//permutations := permute(idSlice, task_map)
	//permutations := gen_all(10, idSlice, task_map)
	//fmt.Println(permutations)

	//fmt.Println()
	//fmt.Println(permutations)
	//fmt.Println()

	maxPers := 5

	var ch chan Pair = make(chan Pair)

	// for _, p := range permutations {
	// 	go calc(task_map, p, maxPers, ch)
	// }

	for i := 0; i < 10; i++ {
		go calc(task_map, generate(idSlice, task_map), maxPers, ch)
	}

	var min_time int = 10000
	var res_comb = make([]int, 0)
	for i := 0; i < 10; i++ {
		pair := <-ch
		if pair.time < min_time {
			min_time = pair.time
			res_comb = pair.comp
		}
	}

	fmt.Println()
	fmt.Println("Минимальное время: ")
	fmt.Println(min_time)
	fmt.Println(res_comb)
	fmt.Println()

	var input string
	fmt.Scanln(&input)
}
*/
