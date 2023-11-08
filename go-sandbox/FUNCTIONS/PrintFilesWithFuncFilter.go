/*
	В примере этого урока мы рассматривали функцию PrintAllFilesWithFilterClosure. В качестве параметра она принимает обязательную строку из пути файла, имя которого выводится на печать.
	На её основе напишите функцию PrintFilesWithFuncFilter(path string, predicate func (string) bool).
	В качестве второго параметра принимается функция, которая проверяет свой аргумент на соответствие определённому условию. Если оно выполняется, то функция возвращает true.
	Для примера может быть передана такая функция:
	// containsDot возвращает все пути, содержащие точки
	func containsDot(s string) bool {
		return strings.Contains(s, ".")
	} 
*/

func PrintAllFilesWithFuncFilter(path string, predicate func(string) bool) {
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
            // печатаем имя элемента, если predicate вернёт true
            if predicate(filename) {
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