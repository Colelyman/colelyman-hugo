+++
title = "Modelling Gravity using the Euler Method in Haskell"
author = ["Cole Lyman"]
date = 2018-09-05T22:54:00-06:00
draft = false
creator = "Emacs 26.1 (Org mode 9.1.13 + ox-hugo)"
+++

## Introduction {#introduction}

In my predictive modelling class, we were given a homework assignment to apply
the [Euler method](https://en.wikipedia.org/wiki/Euler%5Fmethod) for a simple equation and then write a program to
compute the method. I will show what I did to apply the Euler method, and then
how I programmed it using Haskell.


## Our Equation {#our-equation}

For our task we were given the simple equation to run the Euler method on:

\begin{equation}
\frac{dv}{dt} = -g
\end{equation}

where \\(g=9.8 m/s^2\\) (this should look familiar to any physicists). This equation
models the falling of an object (on Earth).


## Applying the Euler Method {#applying-the-euler-method}

To apply the Euler method, one must first have a base case. In this instance,
assume that the initial velocity of the object is 0, thus when \\(t = 0\\), \\(v = 0\\);
in other words \\(v(0) = 0\\).

Now, from my understanding the Euler method is basically a recursive way to
approximate the function. Thus, let \\(h\\) be a given step size and \\(v(t)\\) be the
velocity at a given time \\(t\\), and \\(v\_n\\) be the velocity at a given step \\(n\\). It
follows by Euler's method that:

\begin{equation}
v\_{n+1} = v\_n + hv(t)
\end{equation}

In English, this equation is saying that the velocity at the next time step is
equal to the velocity at the previous time step (\\(v\_n\\)) plus the velocity of the
current time (\\(v(t)\\)) times the step size (\\(h\\)). How do you compute the velocity
at the previous time step (\\(v\_n\\))? Well, it is \\(v\_n = v(t - h)\\), the velocity at
the previous time step.


## Translating Math into Haskell {#translating-math-into-haskell}

Due to Haskell's nature, it isn't hard to translate mathematical equations into
Haskell functions (if you have ever seen a Haskell function, you know what I
mean). Because of the recursive nature, the reasoning behind the Haskell
function is a little bit different than the mathematical intuition presented
earlier, but it is the same principle.

First, let's define \\(g\\) in Haskell:

```haskell
g :: Float
g = 9.8
```

This is pretty straightforward, as defined earlier \\(g = 9.8\\) (the acceleration
of gravity on Earth). If you are unfamiliar with Haskell, the first line
describes the type of \\(g\\) and the second line sets the value of \\(g\\). Now let's
move onto the meat of the algorithm!

Forming the function step by step, the first step is to establish the types that
our function will be dealing with:

```haskell
v :: Float -> Float -> Float
```

Which means that our function `v` will take two parameters (that are of type
`Float`), and return a value that is a `Float`.

Next, let's create our function modelled after the formulae we presented
earlier:

```haskell
v :: Float -> Float -> Float
v t h = v_n + v_t
    where v_n = v (t - h) h
          v_t = h * (-g)
```

this should look very similar to the equation above. However, the acute observer
will notice that this function is not complete because there is no base case for
the recursion to break, thus this function will run forever!!

At the beginning we assumed that our base case was that the initial velocity was
0, thus adding this to our function we get:

```haskell
v :: Float -> Float -> Float
v t h | t > 0 = v_n + v_t
      | othwerwise = 0.0
    where v_n = v (t - h) h
          v_t = h * (-g)
```

In Haskell the `|` are called [guards](https://wiki.haskell.org/Pattern%5Fguard), which essentially translates into a
conditional statement where the first expression that evaluates to true on the
left hand side of the `=` returns what is on the right hand side of the `=`.

For example, if we test our function where \\(t = 1\\) and \\(h = 1\\), then the first
time the function is called the first guard will evaluate to true (because \\(t =
1; 1 > 0 = = True\\)) thus the function is called again (inside `v_n`), but \\(t =
0\\) so the second guard is reached. Ultimately the function will return \\(9.8\\).


## Conclusion {#conclusion}

This was my first encounter using the Euler method to do numerical analysis and
model an equation. I hope that you found this interesting and that you enjoy the
elegance of Haskell as much as I do!
