# Go Cheat Sheet

## Tool Commands

### Cross Compiling
`go build` works great for building on your current OS/architecture. To target a different one:

```bash
GOOS="linux" GOARCH="amd64" go build hello.go
file hello   # check the resulting binary
```

### go doc
Lets you look up documentation for a package, function, method, or variable:

```bash
go doc fmt.Println
```

---

## Data Types

### Primitive Types
`bool`, `string`, `int`, `int8/16/32/64`, `uint`, `uint8/16/32/64`, `uintptr`, `byte`, `rune`, `float32/64`, `complex64/128`

```go
var x = "Hello World"
z := int(42)
```

### Slices & Maps
- **Slice**: like an array but dynamically resizable, more efficient to pass around.
- **Map**: an associative array (key/value pairs), unordered.

```go
var s = make([]string, 0)
var m = make(map[string]string)
s = append(s, "some string")
m["some key"] = "some value"
```

### Pointers
Point to a location in memory.

```go
var count = int(42)
ptr := &count        // & = takes the address
fmt.Println(*ptr)     // * = dereference, reads the value
*ptr = 100            // changes the value in memory
fmt.Println(count)    // 100
```

### Structs
Define new types made up of fields and methods.

```go
type Person struct {
    Name string
    Age  int
}

func (p *Person) SayHello() {
    fmt.Println("Hello, ", p.Name)
}

func main() {
    var guy = new(Person)
    guy.Name = "Dave"
    guy.SayHello()
}
```
> `p` acts like `this/self` in other languages. `new(Person)` creates a new instance (returns a pointer).

### Interfaces
Act like a "contract" — they define which methods a type must implement.

```go
type Friend interface {
    SayHello()
}

func Greet(f Friend) {
    f.SayHello()
}

func main() {
    var guy = new(Person)
    guy.Name = "Dave"
    Greet(guy)
}
```
> Any type that implements `SayHello()` automatically satisfies `Friend`.

---

## Control Structures

### if / else
```go
if x == 1 {
    fmt.Println("X is equal to 1")
} else {
    fmt.Println("X is not equal to 1")
}
```

### switch
```go
switch x {
case "foo":
    fmt.Println("Found foo")
case "bar":
    fmt.Println("Found bar")
default:
    fmt.Println("Default case")
}
```

---

## Concurrency (Goroutines)

The `go` keyword before a call runs it concurrently:

```go
func f() {
    fmt.Println("f function")
}

func main() {
    go f()
    time.Sleep(1 * time.Second)
    fmt.Println("main function")
}
```

---

## Error Handling

Go has **no** try/catch/finally. Instead it uses the built-in `error` interface:

```go
type error interface {
    Error() string
}
```

Custom error type:

```go
type MyError string

func (e MyError) Error() string {
    return string(e)
}
```

---

## Quick Recap

| Concept | Symbol/Keyword | Note |
|---|---|---|
| Address-of | `&x` | gets the memory address |
| Dereference | `*ptr` | reads/sets the value |
| New struct instance | `new(Type)` | returns a pointer |
| Concurrent execution | `go func()` | starts a goroutine |
| Defining a new type | `type X struct{}` / `type X interface{}` | |
