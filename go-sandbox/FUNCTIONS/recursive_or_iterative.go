//recursive factorial

func fact(n int) int {
    if n == 0 {    // терминальная ветка — то есть условие выхода из рекурсии
        return 1
    } else {    // рекурсивная ветка 
        return n * fact(n-1)
    }
}

//iterative factorial

func fact(n int) int {
    res := 1
    for n > 0 {
        res *= n
        n--
    }
    return res
}

//recursive fibonacci

func Fib(n int) int {
    switch {
    case n <= 1:    // терминальная ветка 
        return n
    default:        // рекурсивная ветка
        return Fib(n-1) + Fib(n-2)
    }
}

//iterative fibonacci

func Fib(n int) int {
    a, b := 0, 1
    for n > 0 {
        a, b = b, a+b
        n--
    }
    return a
}