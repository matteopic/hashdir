hashdir
=======

A fast file hashing tool designed to find duplicate files by performing checksums and indexing them.
Brought to you with love by matteopic and friends, developed in Go.

This tool traverses the directory tree, computing the `xxhash` for each file and tracking the results in a text file along with the file size and path.
The index file is utilized by the `hashdir stats` command to identify duplicate files, even if they have different names.
The final report highlights the largest duplicates found.


### Quickstart ###

    go install https://github.com/matteopic/hashdir/

    # Linux
    hashdir /usr
    hashdir stats

    # Windows
    hashdir C:\
    hashdir stats
