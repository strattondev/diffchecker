# DiffChecker

Api access for [diffchecker.com](https://www.diffchecker.com/).

### Install
```
go get github.com/strattonw/diffchecker
```

### Usage
```go
diffcheckerurl, err := DiffChecker{"email", "password"}.Upload("1", "2", "Test")

if err != nil {
	panic(err)
}

fmt.Println(diffcheckerurl)
```
Requires [diffchecker account](https://www.diffchecker.com/signup).

