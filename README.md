[![CircleCI](https://circleci.com/gh/podhmo/handy.svg?style=svg)](https://circleci.com/gh/podhmo/handy)


# handy

## motivation

Sometimes, simple is too strict, and easy is too complicated

## concepts

testing library is error check library

## examples

00

```go
// *testing.T
if got := add(10, 20); got != 30 {
	t.Errorf("want "%d", but got %d", got)
}

// handy
handy.Assert(t, handy.Equal(add(10,20)).Expected(30))

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
handy.Assert(t, handy.EqualNoError(Count(x)).Expected(10))

// testify/assert
c, err := Count(x)
require.NoError(t, err)
assert.Exactly(t, c, 10)
```
