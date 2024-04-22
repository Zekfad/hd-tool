# HD Tool

Utility to work with modified Stingray packages.

Main features:

* Collect hashes (Hash DB Target) from packages (file names, types, package
  names).
* Hash DB:
  * Update via string files.
  * Filter by Hash DB Target.
  * Sort de-hashed entries by values (in natural order).
* Compute hash value of a string.
* Search for a file with a type in packages.

> [!Note]
> Requires `GOEXPERIMENT=rangefunc`
