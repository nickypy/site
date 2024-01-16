---
title: Notes on Python Protocols
published: true
date: 2024-01-15T09:30:04.787298-08:00
slug: "/blog/python-protocols"
tags:
  - python
  - mypy
---
Something I’ve run into often is an unannotated Python function that accepts an argument that calls a method:

```python
def foo(f):
    f.bar()
```

I’d like to annotate this function to gain some benefits of type checking. While I *could* go find all the possible types and `Union` them all, that doesn’t quite work for consumers of this API where arbitrary objects could be supplied. Ideally, the type checker should be happy when the argument implements `bar(self) -> None`. 

Protocols in Python (PEP 544) allow for statically checking whether an object implements a specific method (see interfaces in Go, or traits in Rust). Protocols are also referred to as structural subtyping, or static duck typing.

## Introduction

We can define a protocol by subclassing `Protocol` from the `typing` module:

```python
from typing import Protocol

class Foo(Protocol):
    def bar(self) -> None:
        ...
```

Now that it’s defined, we can annotate our function . The example below results in a happy type checker.

```python
def foo(f: Foo) -> None:
    f.bar()

class A:
    def bar(self) -> None:
        pass

foo(A())
```

Let's see what happens when we use invalid object that does not implement `.bar()`:

```python
class B:
    pass

foo(B())
```

Oof. The type checker gave us a piece of its mind.

```bash
error: Argument 1 to "foo" has incompatible type "B"; expected "Foo"  [arg-type]
Found 1 error in 1 file (checked 1 source file)
```

## Checking a Protocol at runtime

While protocols allow for ahead-of-time type checking, they do not immediately have support for checking types at runtime. To use `isinstance()` with a protocol, we need to annotate the class with `@runtime_checkable`.

```python
from typing import Protocol, runtime_checkable

@runtime_checkable
class Foo(Protocol):
    def bar(self) -> None:
        ...

def foo(f: Foo) -> None:
		assert isinstance(f, Foo)
    f.bar()
```

## Python has some Protocols built-in

Python’s `typing` module contains a few Protocols out of the box, such as `Sized` and `Iterable`. Let’s try those out.

```python
from typing import Iterable, Sized

# The `Sized` protocol implements __len__
def get_size(s: Sized) -> int:
    return len(s)

get_size([1, 2, 3])

# The `Iterable` protocol implements __iter__
def iterate(i: Iterable) -> None:
    for _ in i: 
				pass

iterate([1, 2, 3])
```

## Multiple Protocols

We can define a Protocol that implements multiple protocols via multiple inheritance. Note: we need to subclass `Protocol` , even though the parent classes also subclass `Protocol`.

```python
class SizedIterable(Iterable, Sized, Protocol):
    ...

def do_something(_: SizedIterable) -> None:
    ...

do_something([1, 2, 3])
```

## Conclusion

When Protocols are used correctly, you should see the following `mypy` output:

```python
$ mypy your_file.py
Success: no issues found in 1 source file
```

Protocols provide a convenient way to type-check whether the supplied arguments are valid. They are useful for scenarios where accepting concrete types is either too verbose or not possible. The [original PEP](https://peps.python.org/pep-0544/) contains even more ways to use protocols, motivations behind the protocol, as well as reasons for implementing things the way they are.
