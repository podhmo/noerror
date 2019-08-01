[![CircleCI](https://circleci.com/gh/podhmo/handy.svg?style=svg)](https://circleci.com/gh/podhmo/handy)

# handy

## Motivation

Sometimes, simple is too strict, and easy is too complicated

## Concepts

- testing library is error check library
- zero dependencies

## Examples

00

```go
// *testing.T
if got := add(10, 20); got != 30 {
	t.Errorf("want "%d", but got %d", got)
}

// handy
handy.Assert(t, handy.Equal(30).Actual(add(10,20)))

// testify/assert
assert.Exactly(t, 30, add(10,20))
```

01 (TODO)

```go
// *testing.T
c, err := Count(x)
if err != nil {
	t.Fatalf("error is occured %v", err)
}
if c != 10 {
	t.Fatalf("want "%d", but got %d", 10, got)
}

// handy
handy.Assert(t, handy.Equal(10).ActualWithNoError(Count(x))

// testify/assert
c, err := Count(x)
require.NoError(t, err)
assert.Exactly(t, c, 10)
```
