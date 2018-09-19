+++
title = "Hamming Distance One-liner in Python"
author = ["Cole Lyman"]
date = 2018-09-19T08:29:00-06:00
draft = false
creator = "Emacs 26.1 (Org mode 9.1.13 + ox-hugo)"
+++

The hamming distance of strings \\(a\\) and \\(b\\) is defined as the number of
character mismatches between \\(a\\) and \\(b\\). The hamming distance can be calculated
in a fairly concise single line using Python.


## The Code {#the-code}

```python
def hamming(a, b):
    return len([i for i in filter(lambda x: x[0] != x[1], zip(a, b))])
```


## The Explanation {#the-explanation}

If you aren't familiar with many of the shortcuts that Python provides, this
line may seem quite cryptic. To best explain how this function works, I will
expand it to the equivalent multi-line version, and then go over each part of
the multi-line version.

```python
def hamming(a, b):
    zipped_strings = zip(a, b)

    mismatched_chars = []

    # equivalent of the [...]
    for i in zipped_strings:
        # equivalent of the filter(lambda x: x[0] != x[1], ...)
        if i[0] != i[1]:
            mismatched_chars.append(i)

    return len(mismatched_chars)
```


### Zip {#zip}

`zip` is a useful built-in Python function that takes two lists as its arguments
and returns a list where each element is a tuple, and the first element in the
tuple comes from the first list and the second element in the tuple comes from
the second list. For example:

```python
a = [1, 2, 3]
b = ['a', 'b', 'c']

print(list(zip(a, b)))
```

```text
[(1, 'a'), (2, 'b'), (3, 'c')]
```

In this case, the elements of `a` are _paired up_ with the elements of `b`.


### List comprehension {#list-comprehension}

List comprehensions in Python are a very handy trick to shorten any `for` loop
in Python. They follow the form of `[... for i in iterable]` where the `...` is
replaced by the code that is run on each element of `iterable` (list, generator,
etc.).

An example of a list comprehension is:

```python
evens = [i * 2 for i in range(10) if i % 2 == 0]

print(evens)
```

This code creates a list of all even numbers between 0 and 10, each of which are
multiplied by 2. Notice how the conditional is placed at the end of the expression.


### Lambda functions {#lambda-functions}

Lambda functions may look scary (mostly because it can be hard to recognize
where the parameters come from), but just think of them as functions that don't
have a name. Lambda functions are usually short functions that
perform something straight-forward. In the case of the one-liner, the lambda
function

```python
lambda x: x[0] != x[1]
```

takes one argument `x` (which is of type tuple) and checks if the two elements
of the tuple (`x[0]` and `x[1`) are equal, thereby returning a boolean value.


### Filter {#filter}

The final piece to our one-line puzzle, `filter`. In functional programming
there are patterns and ways to perform certain operations. Three ubiquitous
functions in any language resembling a functional language are `map`, `reduce`,
and `filter`. They are of themselves very simple, yet extremely powerful (if you
want to read more about how to use them in Python, I would recommend [Python
Tips](http://book.pythontips.com/en/latest/map%5Ffilter.html)).

`filter` takes two arguments, first a function that takes an element of an
iterable as input and returns a boolean, and second an iterable. In our one line
case, we have

```python
filter(lambda x: x[0] != x[1], zip(a, b))
```

As described earlier, `lambda` is just a function without a name, and zip is a
list of tuples from two elements of a list. Our `filter` expression returns a
list where the condition `x[0] != x[1]` is `True`, thus giving us a list of
characters that don't match up with one another. When we take the length of this
list we get, by definition, the hamming distance.
