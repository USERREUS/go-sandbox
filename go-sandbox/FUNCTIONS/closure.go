//Классический пример итератора чётных чисел, построенного на замыкании

func Generate(seed int) func() {
    return func() {
        fmt.Println(seed) // замыкание получает внешнюю переменную seed
        seed += 2 // переменная модифицируется
    }
    
}

//Фибоначчи с применением замыкания

func fib() func() int {
    x1, x2 := 0, 1
    // возвращаемая функция замыкает x1, x2
    return func() int {
        x1, x2 = x2, x1+x2
        return x1
    }
}

//Создадим две функции-обёртки, одна из которых будет подсчитывать количество вызовов, а вторая — время исполнения функции

// countCall — функция-обёртка для подсчёта вызовов
func countCall(f func(string)) func(string) {
    cnt := 0
    // получаем имя функции. Подробнее об этом вызове будет рассказано в следующем курсе
    funcname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
    return func(s string) {
        cnt++
        fmt.Printf("Функция %s вызвана %d раз\n", funcname, cnt)
        f(s)
    }
}

// metricTimeCall — функция-обёртка для замера времени
func metricTimeCall(f func(string)) func(string) {
    return func(s string) {
        start := time.Now() // получаем текущее время
        f(s)
        fmt.Println("Execution time: ", time.Now().Sub(start)) // получаем интервал времени как разницу между двумя временными метками
    }
}

func myprint(s string) {
    fmt.Println(s)
}

func main() {

    countedPrint := countCall(myprint)
    countedPrint("Hello world")
    countedPrint("Hi")

    // обратите внимание, что мы оборачиваем уже обёрнутую функцию, а значение счётчика вызовов при этом сохраняется
    countAndMetricPrint := metricTimeCall(countedPrint)
    countAndMetricPrint("Hello")
    countAndMetricPrint("World")

}

//Вспомним функцию PrintAllFilesWithFilter, с которой мы недавно работали. Её недостаток в том, что параметр filter передаётся при каждом рекурсивном вызове. От этого можно избавиться, используя анонимную функцию в качестве замыкания

func PrintAllFilesWithFilterClosure(path string, filter string) {
    // создаём переменную, содержащую функцию обхода
    // мы создаём её заранее, а не через оператор :=, чтобы замыкание могло сослаться на него
    var walk func(string)
    walk = func(path string) {
        // получаем список всех элементов в папке (и файлов, и директорий)
        files, err := os.ReadDir(path)
        if err != nil {
            fmt.Println("unable to get list of files", err)
            return
        }
        //  проходим по списку
        for _, f := range files {
            // получаем имя элемента
            // filepath.Join — функция, которая собирает путь к элементу с разделителями
            filename := filepath.Join(path, f.Name())
            // печатаем имя элемента, если путь к нему содержит filter, который получим из внешнего контекста
            if strings.Contains(filename, filter) {
                fmt.Println(filename)
            }
            // если элемент — директория, то вызываем для него рекурсивно ту же функцию
            if f.IsDir() {
                walk(filename)
            }
        }
    }
    // теперь вызовем функцию walk
    walk(path)
}

