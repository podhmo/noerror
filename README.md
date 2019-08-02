[![CircleCI](https://circleci.com/gh/podhmo/noerr.svg?style=svg)](https://circleci.com/gh/podhmo/noerr)

# noerr

## Motivation

Sometimes, simple is too strict, and easy is too complicated

## Concepts

- testing library is error check library
- zero dependencies

## How to use

Two types of functions are defined.

- compare functions
- assertion functions (or error check functions)

Usually, use these functions as follows.

```go
// t is *testing.T
// want is expected value
// got is actual value

<Assertion>(t, <Compare>(want).Actual(got))
```

e.g.

```go
noerr.Must(t, noerr.Equal(30).Actual(add(10,20)))
noerr.Should(t, noerr.NotEqual(30).Actual(add(10,20)))
```

### Assertion functions

Assertion functions are just error check functions. There are only 2 functions.

- `Must(*testing.T, error)`
- `Should(*testing.T, error)`

`Must()`, if error is occured, calling `*testing.T.Fatal()`. And using `Should()`, `t.Error()` is called when error is occured.

### Compare functions

6 functions are defined.

compare by `x == y`.

- `Equal(expected interface{}) *Handy`
- `NotEqual(expected interface{}) *Handy`

compare by `reflect.DeepEqual(x, y)`.

- `DeepEqual(expected interface{}) *Handy`
- `NotDeepEqual(expected interface{}) *Handy`

compare by `reflect.DeepEqual(normalize(x), normalize(y))`, in default, normalize function is something like `json.Unmarshal(json.Marshal())`.

- `JSONEqual(expected interface{}) *Handy`
- `NotJSONEqual(expected interface{}) *Handy`

## Examples

### normaly usecase.

Normaly, we use `Must()` and `Should()` with compare functions (such as `Equal()`).

```go
// noerr
noerr.Must(t, noerr.Equal(30).Actual(add(10,20)))
noerr.Should(t, noerr.Equal(30).Actual(add(10,20)))


// *testing.T 's example (simple but tiresome)
if got := add(10, 20); got != 30 {
	t.Fatalf("expected "%d", but actual %d", got)
}
if got := add(10, 20); got != 30 {
	t.Errorf("expected "%d", but actual %d", got)
}


// testify/assert 's example (easy but too have methods, so confusing)
require.Exactly(t, 30, add(10,20))
assert.Exactly(t, 30, add(10,20))
```

### One more thing

One more thing, using with `ActualWithNoError()`, support the function returning multiple values, something like `func (<value>) (<value>, error)`. 

If error is occured, computation is stopped and `t.Fatal()` is called.

```go
// Count is something like `func(int) (int, error)`.
//
// e.g.
// func Count(xs []<value>) (int, error) {
// .. do something
// }


// noerr
noerr.Should(t, noerr.Equal(10).ActualWithNoError(Count(x))


// *testing.T
c, err := Count(x)
if err != nil {
	t.Fatalf("error is occured %v", err)
}
if c != 10 {
	t.Fatalf("want "%d", but got %d", 10, got)
}

// testify/assert
c, err := Count(x)
require.NoError(t, err)
assert.Exactly(t, c, 10)
```

## Output if test is failed

code

```go
noerr.Should(t, noerr.Equal(10).Actual(1+1))
noerr.Should(t, noerr.Equal(10).Actual(1+1).Describe("1+1"))

// Count() returns error
Count := func() (int, error) { return 0, fmt.Errorf("*ERROR*") }
noerr.Should(t, noerr.Equal(0).ActualWithNoError(Count()))
```

output

```
Equal, expected 10, but actual 2
1+1, expected 10, but actual 2

unexpected error, *ERROR*
```

custom output examples, [TODO](https://github.com/podhmo/noerr/blob/ab7573214f6da953fef1a17692ff9abf4d09686c/log_test.go#L59)
