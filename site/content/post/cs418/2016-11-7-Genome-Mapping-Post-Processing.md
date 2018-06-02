---
categories: cs418
date: "2016-11-07T00:00:00Z"
tags: school, genome mapping
title: Genome Mapping Post Processing
---

So you have mapped your reads to the reference genome, but what comes next? 
How can you tell how many reads were aligned, where they were aligned, and actually see what your mapping algorithm did?
This post will show you how you can analyze the results of your genome mapper by using [samtools](http://www.htslib.org/) and [IGV](http://software.broadinstitute.org/software/igv/).

# samtools

## SAM File Format

Most modern genome aligners will output the aligned reads in the [SAM file format](http://genome.sph.umich.edu/wiki/SAM).
If you have a lot of time on your hands you can read through this SAM file and see where the reads are mapped to and read the information for each read (which probably adds up to millions of lines).
A'int nobody got time for that.
Instead, we are going to have `samtools` do the work for us.

## Installation

1. [Download samtools](https://github.com/samtools/samtools/releases/download/1.3.1/samtools-1.3.1.tar.bz2).
2. Go to the directory that you downloaded it and unzip the file
  * `$ cd ~/Downloads`
  * `$ tar jxvf samtools-1.3.1.tar.bz2`
3. Go to the upzipped directory and prepare for install
  * `$ cd samtools-1.3.1`
  * `$ ./configure`
    * **Note:** You may get a library dependency error if you don't have the development files for the ncurses library installed. If you do get this error, install either `libncurses5-dev` or `ncurses-devel`.
4. If the `./configure` command executed without errors, install samtools
  * `$ sudo make install`
  * **Note:** If you do not have super-user privileges then you can run `$ make install --prefix=<directory of your choice>`
5. Test if the install worked by running `samtools`
  * `$ samtools`

If the installation was successful then you should see a usage prompt for `samtools`, otherwise your installation failed (or you aren't using the correct path to the `samtools` binary).

## BAM File Format

The [BAM file format](http://genome.sph.umich.edu/wiki/BAM) is the binary cousin to the SAM file format, it is the same thing, just in binary form (this means that you can't directly view a BAM file in a text editor like `vim` like you could a SAM tile).

## Converting SAM to BAM

We will now use `samtools` to convert our SAM file to a BAM file.
The SAM file that we are going to use is found in the `examples` directory of the `samtools-1.3.1` directory, the file is `toy.sam`.
`samtools` is broken up into subprograms that perform different functions.
In order to convert the SAM file we use `samtools view`, like so: `$ samtools view -b -S examples/toy.sam -o examples/toy.bam`.

## Sorting BAM

Next, we want to sort the reads in the BAM file.
This is a necessary step because many visualization programs cannot handle unsorted BAM files.
To sort the BAM file, run: `$ samtools sort examples/toy.bam -o examples/toy.sorted.bam`.

## Indexing BAM

We then have to index the sorted BAM file, which is also a necessary step for visualization tools.
Run `$ samtools index examples/toy.sorted.bam`.

## Viewing BAM Using `samtools tview`

Now we can finally view the aligned reads in our terminal!
Run `$ samtools tview examples/toy.sorted.bam examples/toy.fa` and you will suddenly see four reads aligned to an extremely small genome. 
You should see something along these lines:

    1         11              21        31         41        51
    AGCATGTTAGATAA****GATA**GCTGTGCTAGTAGGCAG*TCAGCGCCATNNNNNNNN
      ........    ....  ......K.K......K. ..........
      ........AGAG....***...      ,,,,,    ,,,,,,,,,
        ......GG**....AA
        ..C...**** ...**...>>>>>>>>>>>>>>T.....```

The numbers at the top signify the index of the genome, and the first line of characters represents the reference genome itself.
The third line is the (consensus sequence)[https://en.wikipedia.org/wiki/Consensus_sequence] (the astute student will notice that in the consensus sequence there are three `K`'s, [look here](http://www.chick.manchester.ac.uk/SiteSeer/IUPAC_codes.html) to discover what `K` represents).
Each line under the consensus sequence is a read. 

You may be wondering why most of the reads are made of `.`, well that is because they match the reference genome.
There are many different settings that you can play with in `samtools tview`, to view all of the settings type `?` and a help menu will come appear.

# Integrative Genome Viewer (IGV)

If you want more flexibility and a more robust way of viewing your aligned reads you can use IGV.
It has a GUI, which makes things nice sometimes.

## Installation

[Download](http://software.broadinstitute.org/software/igv/download) the appropriate files according to the system that you are running. 
IGV is written in Java (using Java 7, not sure if Java 8/9 will work, but it should), so you need to make sure that you have Java installed on your computer.
Once you have Java installed and IGV downloaded, go ahead and unzip the downloaded file, if needed.
Then if you are in an Unix-like OS (Mac or Linux) you can run `$ ./IGV_2.3.88/igv.sh` to open up the IGV GUI.

## Loading the Reference Genome

Once the program is open, click the `Genomes` button at the top of the window, then select `Load Genome from File...`.
Select the file `samtools-1.3.1/examples/toy.fa`. 

## Loading the Reads

After the reference genome is loaded we can load in the reads by clicking the `File` button at the top of the window, then select `Load from File...`.
Select the file `samtools-1.3.1/examples/toy.sorted.bam`.

## Seeing the Reads

You probably can't see any changes in the view, that's ok.
In the second row click on the dropdown arrow that says `All` and select either `ref` or `ref2` and then you will be able to see what we saw in `samtools tview`, except it is way easier to figure out what everything means!
