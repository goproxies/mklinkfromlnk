# Do

On msys2 when `tar` huge file `aosp-last.tar` to windows all `link files` were changed to `*.lnk` format.This tool mainly execute `mklink` command to create soft link(symlink) instead.

# Execution flow

- admin right neaded(which requested by `nac.manifest`)
- parse `*.lnk` file(s) under root dirs first,and then each subdir of them
    - make symlink for every `*.lnk` by executing `mklink`
    - delete *.lnk file(s) if `-e` switch on


# Build
- change current dir to project dir
- `rsrc -manifest nac.manifest -o nac.syso` if nac.manifest changed
- `go install -v -gcflags "-N -l" ./...`

# Release
- latest-release ![v1.0.0-amd64](https://github.com/goproxies/mklinkfromlnk/releases/download/v1.0.0-amd64/mklinkfromlnk.exe)

# Usage:
```shell
Usage: mklinkfromlnk [-h] [-e] [-r] [-s] [-printerror] [-printstack [-printerror]] [-v [-s -printerror]] [-oledebug [-printstack -printerror]] [-vv [-v -printstack -oledebug]] [-c maximum-numbers-of-coroutines] [-r use-relative-path] [-d directory] [directory...]
e-path] [-d directory...]

Options:
  -c int
        minimum 20,maximum numbers of coroutines  (default 50)
  -d string
        which directory to search (default ".")
  -e    delete *.lnk
  -h    print help
  -oledebug
        show ole errors
  -printerror
        print errors
  -printstack
        print stack when some errors occured
  -r    create with relative path
  -s    show parsing dir
  -skip string
        skip directory pattern
  -v    show infos verbosely
  -vv
        show infos more verbosely
```
# Example:

The follow command create `symlink` with relative path, delete the `.lnk` file and skip all objects dirs `objects\..` of git repo
- mklinkfromlnk -r -e -skip "objects\\\\[0-9a-z]{2}$"
