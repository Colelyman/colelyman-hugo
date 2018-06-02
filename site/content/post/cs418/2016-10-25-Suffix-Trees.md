---
categories: cs418
date: "2016-10-25T00:00:00Z"
tags: data structures
title: Suffix Trees
---

# Suffix Trees

## What is a Suffix Tree?

A suffix tree is a data sructure that contains each suffix of a string, a suffix is defined as the end part of a string. 
For example, the suffixes for the string banana are:

* a
* na
* ana
* nana
* anana
* banana

Notice that the entire string is also considered a suffix. 
We can also include the empty string as a suffix, so in order to include the empty string we need to append a character that isn't in the alphabet to the string.
We will assume that the character `$` is not in the alphabet.
Our string is now banana$, and the suffixes are:

* $
* a$
* na$
* ana$
* nana$
* anana$
* banana$

## How to Construct a Suffix Tree

The easiest way to understand how to construct a suffix tree is to first construct a suffix trie, then collapse nodes to convert the trie to a tree.
Here is the graphical representation of the suffix trie:
![suffix trie](/public/cs418/suffixTrie.png)

Now we collapse the nodes that only have one child, and the resulting suffix tree is this:
![suffix tree](/public/cs418/suffixTree.png)

## How to Use a Suffix Tree

Observe that each suffix is represented in this structure. 
This means that we can efficiently search for a suffix (or substring) by traversing the suffix tree.
For example if we were to search for the substring `ana` in the suffix tree for `banana$`, this would be the traversal path:

```0 -> 2 -> 3```

*Note:* We can stop at any node because we are searching for a substring rather than a suffix, if we were searching for a suffix, the search string would have been `ana$`.

## Suffix Tree Construction

How would you create a naive algorithm to construct a suffix tree?
You can easily create a simple algorithm to construct a suffix tree in _O(n<sup>2</sup>)_, but [Esko Ukkonen](https://www.cs.helsinki.fi/u/ukkonen/) discovered a way to [construct a suffix tree in _O(n)_ (linear time)](http://stackoverflow.com/questions/9452701/ukkonens-suffix-tree-algorithm-in-plain-english/9513423#9513423).
