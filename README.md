# macro
Our macro package aggregates multiple internal little tools, which we want to try out.
Some tools may have been proofed to be useful, for some we still collect usage reports and others are just here for archive or compatibility purposes.

## FAQ

### How to execute?

Change the working directory to your go module and call the macro expander as follows:

```bash
go run github.com/worldiety/macro/cmd/expand@latest
```

Note, that go caches github aggressively, so ensure that you have the latest version, if new commits have arrived today:

```bash
GOPROXY=direct go run github.com/worldiety/macro/cmd/expand@latest
```

### There are so many error logs
Sorry, but the implementation is not (yet) complete.
However, the result is often still acceptable.
Try to switch to supported basic types and/or introduce distinct named types.

### Where are the unit tests?
They wait to get written by you.

## tagged union
In the following FAQ we document our discussions about the tagged union macro.

### Example

Consider some source code like

```go
// A Component is a sum type or tagged union.
// Actually, we can generate different flavors, so that Go makes fun for modelling business stuff.
//
// #[go.TaggedUnion "json":"intern", "tag":"type"]
type _Component interface {
	Button | TextField | Text | Chapter | xcompo.RichText | xcompo.Icon | string | []string | []Text
}

```

Note the macro invocation.
Execute the generator as follows to let a `component.gen.go` file be generated aside the file which contains the `_Component` definition:

Take a look at the result at https://github.com/worldiety/macro/blob/main/testdata/example/domain/component.gen.go

### Config options
Currently, only the following macro invocation is possible:

```rust
#[go.TaggedUnion "json":"intern", "tag":"type", "names":[]]
```

which is the same as ommitting the options and resorting to the default settings.

Note, that the names attribute may contain a list of alternate serialization names.
This is important for refactoring type names or just supporting names from external systems.
The order and length of the names array must exactly match the union types order.

```rust
#[go.TaggedUnion]
```

### Why choice types in Go, there are already interfaces?
A choice type can be expressed in many ways.
For us, the developer ergonomic is the most important argument.
A developer has to express a domain model in a readable and understandable way.
An evolution of the domain and therefore the code, must not result in unwanted side effects, thus we want as much as possible support from the type system.

However, it is currently nearly impossible to express a closed (polymorphic) sum type using interfaces in Go.
See also https://github.com/golang/go/issues/57644 for details.
Probably the most popular approach is to introduce package-private methods and an according marker interface.
However, there are edge cases when embedding such types and this also only works for types in the same package.
There are also situations, where a generic instantiation may be (mis-)used to express sum type facts, but these fall immediately short in polymorphic use cases.

### Why not just a linter?
We have considered this, but this means it must run more or less all the time and there is no IDE support like autocompletion for it.
In contrast to that, once generated, the sum type can be used by any other subdomain or as a supporting domain in a type safe way.
A linter does also not help with serialization.
If interfaces as sum types become a thing in Go, we will reconsider a linter.

### Why not just Rust?
Rust is likely also a good choice, however, it may introduce a lot of unwanted technical burden into a domain model.
The ownership model enforces domain facts which are either undefined or sometimes even not sound from the domain perspective.
By definition, the Rust approach excludes a lot of valid solutions, for the sake of limitations imposed by the programming model.
Therefore, a GC oriented language like Go can also be a good fit, just with another tradeof regarding the type system.

### Why a tagged union, isn't that the same as an interface?
Sort of, however, the generated container has a clear discriminator to distinguish each case.
It is not possible to generally solve a static exhaustive match using the Go interface semantics without introducing further restrictions, see also https://blog.merovius.de/posts/2022-05-16-calculating-type-sets/ for more details.
Therefore, we introduce the restriction of the discriminator or enumerator for our type sets.
We do not care if the types within a tagged union are polymorphic interchangeable, they just need to have a distinct name.

There are a lot of other approaches with different trade offs which we want to try out it in the future.

### Why all these AsX() (x,ok) methods?
First of all, a (non package local) type can be part of different choice types within the same package.
If it is a function, the parameter must be polymorphic to all possible choice types.
So, the only natural criteria to identify this, is a method which tells us exactly that fact.
We prefixed it with 'As' to lower the risk of method set collisions, especially with typical ordinary getters which usually do not impose a tuple return.
It is used exactly in the way, as an ordinary type switch could be used, but statically proofed.
These methods enables also a lot of expressive polymorphic cases, where a receiver can just define its accessor interface and can work with any tagged union which may provide it as an element.


### Why a Switch method?
Due to the limited type system in Go (no useful monads), it is idiomatic to use func closures for causing side effects.
We hope, that an alternate pipeline notation will rise in the future and help in reading nested function calls.

### Why function MatchX?
This introduces another level of type safety, which cannot be statically proofed using side effect closures.
Using this functional approach, a developer can be sure to map each case exhaustively and not forget any case by accident.
Any change will more likely result in a compiler error instead of just accidently processing a zero value.

### Why JSON Marshal/Unmarshal?
Unmarshalling interface types in Go has never been a supported feature and always requires a lot of boilerplate code.
Our implementations capture the most common use cases by default.
