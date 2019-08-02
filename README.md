[![CircleCI](https://circleci.com/gh/podhmo/handy.svg?style=svg)](https://circleci.com/gh/podhmo/handy)

# handy

## Motivation

Sometimes, simple is too strict, and easy is too complicated

## Concepts

- testing library is error check library
- zero dependencies

## How to use

Two group functions are defined.

- compare functions
- assertion functions (or error check functions)

Usually, use these functions, like a below.

```go
// t is *testing.T

<Assertion>(t, <Compare>(want).Actual(got))
```

e.g.

```go
handy.Must(t, handy.Equal(30).Actual(add(10,20)))
handy.Should(t, handy.NotEqual(30).Actual(add(10,20)))
```

### Assertion functions

Assertion functions are just error check functions. There are only 2 functions.

- `Must(*testing.T, error)` -- if error is occured, calling t.Fatal()
- `Should(*testing.T, error)` -- if error is occured, calling t.Error()

### Compare functions

6 functions are defined.

Below functions are the functions that compare by `x == y`.

- `Equal(expected interface{}) *Handy`
- `NotEqual(expected interface{}) *Handy`

Below functions are the functions that compare by `reflect.DeepEqual(x, y)`.

- `DeepEqual(expected interface{}) *Handy`
- `NotDeepEqual(expected interface{}) *Handy`

Below functions are the functions that compare by `reflect.DeepEqual(normalize(x), normalize(y))`, in default, normalize function is something like `json.Unmarshal(json.Marshal())`.

- `JSONEqual(expected interface{}) *Handy`
- `NotJSONEqual(expected interface{}) *Handy`

## Examples

### normaly usecase.

Normaly, we use `Must()` and `Should()` with compare functions (such as `Equal()`).

```go
// handy
handy.Must(t, handy.Equal(30).Actual(add(10,20)))
handy.Should(t, handy.Equal(30).Actual(add(10,20)))


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
// Count is something like func(int) (int, error)
// e.g.
// func Count(xs []<value>) (int, error) {
// .. do something
// }


// handy
handy.Should(t, handy.Equal(10).ActualWithNoError(Count(x))


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
