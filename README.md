## What are canonical tests?

That's when instead of comparing the expected and actual values in code:
```go
assert.Equal(t, expected, actual)
```

You instead assert that the actual value is equal to a known canonical value stored in a file `canonical.json`:
```go
canonical.Assert(t, actual)
```

This approach is especially useful when comparing large values.

## Okay, that sounds great, how do I use this library?
First, download the library with `go get github.com/bennydictor/canonical`.

At the top of your test file, add the line
```go
//go:generate go test . -ldflags "-X 'github.com/bennydictor/canonical.Canonize=true'"
```

Run `go generate ./your/package/with/tests`. This will generate the `canonical.json` file.
Look at it and see if the values there make sense, and if not, fix your code.
Add `canonical.json` to your version control system.

`Assert` and `Require` will now compare asserted values with the canonical values saved in `canonical.json`.

Also, there's an [../example](example).

## How do I assert multiple values?
There must be at most one call to `Assert` or `Require` per test. If you want to assert multiple values,
you will have to pass in multiple values:
```go
canonical.Assert(t, foo, bar, baz)
```

Or like that:
```go
canonical.Assert(t, map[string]interface{}{
	"foo": foo,
	"bar": bar,
	"baz": baz,
})
```

## It didn't assert anything! There's just a `{}` in `canonical.json`
As you may have noticed, the values are converted to json in order to be stored and compared.
If your value has private fields anywhere in it, those private fields will not be included in
the json. That's just how Go works, I can't do anything about it.

Try changing the fields to be public. If that's impossible or you don't want to,
then you'll have to assert them one by one (e.g. `canonical.Assert(t, result.foo, result.bar)`).

I made a wrapper for errors tho: `canonical.Assert(canonical.Error(err))` (works with `nil` errors).

## I changed the code and now my tests are failing! What do I do?
Look at the test diff. If you agree with it, then run `go generate ./your/package/with/tests`.
If you don't agree with the diff, fix your code.

## `go generate` overwrote `canonical.json` when i didn't want it to! What do I do?
Get the version you want from your version control system.
If you don't use a version control system, sorry, I can't help you.
