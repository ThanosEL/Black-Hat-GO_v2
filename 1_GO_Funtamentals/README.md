## Using Common GO Tool Commands
### Cross Compilining
Using go build works great for running a binary on your current system or one of identical architecture.
But we can set adn choose different architecture and OS.

```shell
GOOS="linux" GOARCH="amd64" go build hello.go
```

Then check with:
```shell
file hello
```

### The go doc Command
The go doc command lets you interrogate documentation about a package, function, method, or variable.

```shell
go doc fmt.Println
```

## Understading GO Syntax
### Data Types
#### Primitive Data Types
The primitive types include bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, byte, rune, float32, float64, complex64, and complex128.

```go
var x = "Hello World"
z := int(42)
```

#### Slice and Maps
Slices are like arrays that you can dynamically resize and pass to functions more effi-cient ly. Maps are associative arrays, unordered lists of key/value pairs that allow you to efficiently and quickly look up values for a unique key.

```go
var s = make([]string, 0)
var m = make(map[string]string)
s = append(s, "some string")
m["some key"] = "some value"
```

#### Pointers, Structs and Interfaces
A pointer points to a particular area in memory and allows you to retrieve the value stored there.
```go
var count = int(42)
ptr := &count
fmt.Println(*ptr)
*ptr = 100
fmt.Println(count)
```

1. The code defines an integer, count
2. and then creates a pointerby using the & operator. This returns the address of the count variable
3. You dereference the variable while making a call to fmt.Println()
4. You then use the * operator to assign a new value to the memory location pointed to by ptr
5. The assignment changes the value of that variable, which you confirm by printing it to the screen


You use the struct type to define new data types by specifying the type’s associated fields and methods
```go
type Person struct {
    Name    string
    Age     int
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

1. The code uses the type keyword to define a new struct containing two fields: a string named Name and an int named Age.
2. You define a method, SayHello(), on the Person type assigned to variable p
3. The method prints a greeting message to stdout by looking at the struct, p. Think of p as a reference to self or thisin other languages.
4. his function uses the new keyword to initialize a new Person.
5. It assigns the name Dave to the person and then tells the person to SayHello()


You can think of Go’s interface type as a blueprint or a contract.
```go
type Friend interface {
    SayHello()
}
```

1. In this sample, you’ve defined an interface called Friend that requires one method to be implemented: SayHello()
2. That means that any type that implements the SayHello() method is a Friend.

The following function, Greet(), takes a Friend interface as input and says hello in a Friend-specific way:
```go
func Greet (f Friendv) {
    f.SayHello()
}
```

```go
func main() {
    var guy = new(Person)
    guy.Name = "Dave"    
    Greet(guy)
}
```

### Control Structures
Go’s primary conditional is the if/else structure:
```go
if x == 1 {
    fmt.Println("X is equal to 1")
} else {
    fmt.Println("X is not equal to 1")
}
```

For conditionals involving more than two choices, Go provides a switchstatement. The following is an example:
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

### Concurrency
To create a goroutine, use the go keyword before the call to a method or function you wish to run concurrently:
```go
func f() {
    fmt.Println("f function")
}

func main() {
    go f()
    time.Sleep(1 * time.Second)
    fmt.Println("main function")}
```

### Error Handling
Unlike most other modern programming languages, Go does not include syntax for try/catch/finally error handling. 

Go defines a built-in error type with the following interface declaration:
```go
type error interface {
    Error() string
}
```

custom error you could define and use throughout your code:
```go
type MyError string
func (e MyError) Error() string {
    return string(e)
}
```
