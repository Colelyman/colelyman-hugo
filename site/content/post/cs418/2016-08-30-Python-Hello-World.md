---
categories: cs418
date: "2016-08-30T00:00:00Z"
tags: python, school
title: Python Tutorial CS 418/BIO 365
---

<div align="center">
<h1>Python Tutorial</h1>

<h3>CS 418/BIO 365</h3>
<h3>30 August 2016</h3>
</div>


```python
# This is a comment
```


```python
""" This is a 
    multi-
    line
    comment"""
```


```python
''' This is also a 
    multi-
    line
    comment'''
```


```python
# How do I declare variables?

# String
my_string = 'Hello, CS 418/BIO 365!!'
my_other_string = "Hello, friends."
```


```python
# Numbers
my_number = 418
my_other_number = 365.578
```


```python
# List
my_list = ['this', 'is', 'fun', 56]
# or
my_other_list = []
# or
my_other_list2 = list()

my_other_list.append('bioinformatics')
```


```python
# Dictionary (Map)
my_dict = {'key': 'value', 'other_key': 'other_value'}
# or
my_other_dict = {}
# or
my_other_dict2 = dict()

my_other_dict['my_key'] = 'my_value'
```


```python
# How do I make loops?

# NOTE: Whitespace matters!!

# for loop
for i in my_list:
    print(i)
    
# NOTE: "print(i)" is used in Python 3, "print i" is used in Python 2 
```

    this
    is
    fun
    56



```python
# for loop over range
for i in range(0, 5):
    my_other_list.append(i)
print(my_other_list)
```

    ['bioinformatics', 0, 1, 2, 3, 4]



```python
# while loop
while len(my_list) > 0:
    del my_list[-1]
    print(my_list)
```

    ['this', 'is', 'fun']
    ['this', 'is']
    ['this']
    []



```python
# break 
for val in my_other_list:
    if val is 2:
        break
    else:
        print(val)
```

    bioinformatics
    0
    1



```python
# continue
for val in my_other_list:
    if val is not 2:
        continue
    print(val)
```

    2



```python
# How do I perform file I/O?

# file input
import sys # this imports the library to get command line arguments
with open(sys.argv[1]) as fh: # sys.argv[1] is the first command line argument, sys.argv[2] is the second ... and so on.
    my_first_line = next(fh) # Python 2 & 3
    my_next_line = fh.next() # Python 2 
    
with open(sys.argv[1]) as other_fh:
    # iterate over all lines in file
    for line in other_fh:
        print(line)
```

    ACGTTGCATGTCGCATGATGCATGAGAGCT
    
    4
    



```python
# file output
with open('./my_file.txt', 'w') as writable: ''' NOTE: the second 
    argument ('w') makes this file writable '''
    writable.write('This is my test file\n')
    writable.write(my_first_line)
```


```python
# How do I manipulate strings?

# string slice
my_string = 'banana'
print(my_string[3])
print(my_string[0:5])
print(my_string[1:5])
print(my_string[-1])
print(my_string[0:-2])
```

    a
    banan
    anan
    a
    bana


<div align="center">
<h1>Rosalind Tutorial</h1>
<h4>You can register for the course at: <a href="http://rosalind.info/classes/enroll/e464d00a10/">http://rosalind.info/classes/enroll/e464d00a10/</a></h4>
<h3>Use your NetID as your username!!</h3>
</div>

<div align="center">
<h2>The first problem is found here: <a href="http://rosalind.info/problems/ba1b/?class=322">http://rosalind.info/problems/ba1b/?class=322</a></h2>
</div>


```python
# What do we need to do in pseudo code?

''' -Read in the sequence and the 
    length of the kmer (a kmer is essentially a substring)
    -Break up the sequence into kmers
    -Count each kmer
    -Find the kmer(s) with the highest counts
    -Print out the highest count kmers'''
```


```python
# Read in the sequence and the length of the kmer from a file
import sys
with open(sys.argv[1]) as file:
    seq = next(file).strip() # REMEMBER: in Python 2 use file.next()
    kmer_len = int(next(file).strip())
```


```python
# Break up the sequence into kmers
    for i in range(0, len(seq)):
        kmer = seq[ i : i + kmer_len]
```


```python
# Count each kmer
    counts = {}
    for i in range(0, len(seq)):
        kmer = seq[ i : i + kmer_len]
        if kmer not in counts:
            counts[kmer] = 1
        else:
            counts[kmer] += 1
```


```python
# Find the kmer(s) with the highest counts
    max_count = 0
    max_kmers = []
    for kmer, count in counts.items(): # NOTE: use counts.iteritems() in Python 2
        if count > max_count:
            max_count = count
            max_kmers.clear()
            max_kmers.append(kmer)
        elif count == max_count:
            max_kmers.append(kmer)
```


```python
# Print out the kmers with the highest counts
    print(" ".join(max_kmers))
```
