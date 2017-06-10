Local Cache
===========

In order to speed up certain operations, restic manages a local cache of data.
This document describes the data structures for the local cache with version 1.

Versions
--------

The cache directory is selected according to the `XDG base dir specification
<http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html>`__. It
contains a file named `version`, which contains a single ASCII integer line
that stands for the current version of the cache. If a lower version number is
found the cache is recreated with the current version. If a higher version
number is found the cache is ignored and left as is.

Snapshots and Indexes
---------------------

Snapshot, Data and Index files are cached in the sub-directories ``snapshots``,
``data`` and  ``index``, as read from the repository.

Keys
----

Keys are never cached locally. The sub-directory ``key`` contains an empty file
which has the same name as the key file that was successfully used the last
time to open the repository. On a subsequent run, restic tries this key file
first (if it still exists).
