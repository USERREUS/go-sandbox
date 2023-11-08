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
	storePath  = "../storage/tasks2.json" //путь к файлу с задачами
	maxPers    = 5                        //количество работников
	iterations = 1000                     //количество итераций расчета
	pInfTime   = 10000                    //бесконечное время выполнения
)

//Структура задачи для чтения из файла
type Data struct {
	ID   int   `json:"id"`
	Prev []int `json:"prev"`
	Next []int `json:"next"`
	Time int   `json:"time"`
	Res  int   `json:"res"`
}

//Структура задачи для расчетов
type Task struct {
	Prev []int
	Next []int
	Time int
	Res  int
}

//Результат расчета
type Pair struct {
	time int
	comp []int
}

//Расчет времени выполнения задач
func simulate(taskMap map[int]Task, maxPers int, ch chan Pair) {
	sample := genSample(taskMap)
	tmp := make([]int, len(sample))
	copy(tmp, sample)

	done := make(map[int]bool)     //выполненные задачи
	inProcess := make(map[int]int) //задачи в процессе выполнения
	time := 0                      //модельное время
	pers := maxPers                //количество свободных работников

	for len(sample) > 0 {
		taskID := sample[0]
		task := taskMap[taskID]

		//Проверка на готовность задачи к выполнению
		taskReady := true
		if len(task.Prev) > 0 {
			for _, i := range task.Prev {
				if !done[i] {
					taskReady = false
					break
				}
			}
		}

		//Если задача готова к выполнению и на нее хватает свободных сотрудников
		if taskReady && pers >= task.Res {
			inProcess[taskID] = task.Time
			pers -= task.Res
			sample = sample[1:]
		} else {
			//Приращение модельного времени на минимальное из времен выполняюшихся задач
			minTime := pInfTime
			for _, t := range inProcess {
				if t < minTime {
					minTime = t
				}
			}
			time += minTime

			//Перемещение задач из выполняющихся в выполненные
			for item, t := range inProcess {
				if t == minTime {
					pers += taskMap[item].Res
					done[item] = true
					delete(inProcess, item)
				} else {
					inProcess[item] -= minTime
				}
			}
		}
	}

	//Приращение времени на максимум из выполняющихся задач
	maxTime := 0
	for _, t := range inProcess {
		if t > maxTime {
			maxTime = t
		}
	}
	time += maxTime

	//Сохранение результата в канал
	ch <- Pair{time: time, comp: tmp}
}

//Чтение файла в структуру данных
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

//Перемешивание элементов очереди в случайном порядке
func shflQueue(q *[]int) {
	if len(*q) <= 1 {
		return
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(2)
	if r == 1 {
		rand.Shuffle(len(*q), func(i, j int) {
			(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
		})
	}
}

//Генерация последовательности задач
func genSample(taskMap map[int]Task) (res []int) {
	var q []int                          //очередь задач
	used := make([]bool, len(taskMap)+1) //признак выполнения задачи
	for i := range used {
		used[i] = false
	}

	//Заполнение очереди задачами без предков
	for id, item := range taskMap {
		if len(item.Prev) == 0 {
			q = append(q, id)
		}
	}

	for len(q) > 0 {
		//Перемешивание, извлечение и формирование результата
		shflQueue(&q)
		cur := q[0]
		used[cur] = true
		q = q[1:]
		res = append(res, cur)

		//Добавление новых задач в очередь
		for _, id := range taskMap[cur].Next {
			flag := false
			if !used[id] {
				for _, pr := range taskMap[id].Prev {
					if !used[pr] {
						flag = true
						break
					}
				}
				if flag {
					continue
				}
			}
			q = append(q, id)
		}
	}
	return
}

//Проверка на валидность последовательности задач
func isValid(p []int, task_map map[int]Task) bool {
	used := make(map[int]bool)
	valid := true
	for _, i := range p {
		for _, j := range task_map[i].Next {
			if used[j] {
				valid = false
			}
		}
		used[i] = true
	}
	return valid
}

func main() {
	//Парсинг JSON файла с задачами
	data, err := jsonParse(storePath)
	if err != nil {
		log.Fatalf("Не удалось прочитать файл: %v", err)
	}

	//Вывод прочитанных данных
	fmt.Println("Задачи:")
	for _, rec := range data {
		fmt.Printf(
			"ID : %d; Prev : %v; Next : %v; Time : %d; Res : %d",
			rec.ID,
			rec.Prev,
			rec.Next,
			rec.Time,
			rec.Res,
		)
		fmt.Println()
	}
	fmt.Println()

	//Создание словаря идентификатор - задача
	taskMap := make(map[int]Task)
	for _, rec := range data {
		taskMap[rec.ID] = Task{
			Prev: rec.Prev,
			Next: rec.Next,
			Time: rec.Time,
			Res:  rec.Res,
		}
	}

	//Создание канала для параллельной обработки данных
	ch := make(chan Pair)

	//Запуск горутин с расчетной функцией
	for i := 0; i < iterations; i++ {
		go simulate(taskMap, maxPers, ch)
	}

	//Вычисление минимального времени из рассчитанных
	min_time := pInfTime
	var res_comb []int
	for i := 0; i < iterations; i++ {
		pair := <-ch
		if pair.time < min_time {
			min_time = pair.time
			res_comb = pair.comp
		}
	}

	fmt.Printf("Минимальное время = %d. С %d работниками на наборе задач: ", min_time, maxPers)
	fmt.Println(res_comb)
	fmt.Println()

	fmt.Println("Нажмите Enter...")
	var input string
	fmt.Scanln(&input)
}
