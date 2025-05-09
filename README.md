> **Notice:** This repository is an **independent fork** of
> [cheat/cheat](https://github.com/cheat/cheat). The code that has been
> modified or added by me is licensed under [LICENSE](./LICENSE).

# Cheat

`cheat` allows you to create and view interactive cheatsheets on the
command-line. It was designed to help remind \*nix system administrators of
options for commands that they use frequently, but not frequently enough to
remember.

![The obligatory xkcd](http://imgs.xkcd.com/comics/tar.png 'The obligatory xkcd')

Use `cheat` with [cheatsheets](https://github.com/cheat/cheatsheets).

## Example

The next time you're forced to disarm a nuclear weapon without consulting
Google, you may run:

```sh
cheat view tar
```

You will be presented with a cheatsheet resembling the following:

```sh
# To extract an uncompressed archive:
tar -xvf '/path/to/foo.tar'

# To extract a .gz archive:
tar -xzvf '/path/to/foo.tgz'

# To create a .gz archive:
tar -czvf '/path/to/foo.tgz' '/path/to/foo/'

# To extract a .bz2 archive:
tar -xjvf '/path/to/foo.tgz'

# To create a .bz2 archive:
tar -cjvf '/path/to/foo.tgz' '/path/to/foo/'
```

## Usage

To view a cheatsheet:

```sh
cheat view tar      # a "top-level" cheatsheet
cheat view foo/bar  # a "nested" cheatsheet
```

To edit a cheatsheet:

```sh
cheat edit tar     # opens the "tar" cheatsheet for editing, or creates it if it does not exist
cheat edit foo/bar # nested cheatsheets are accessed like this
```

To view the configured cheatpaths:

```sh
cheat dirs
```

To list all available cheatsheets:

```sh
cheat ls
```

To list all cheatsheets that are tagged with "networking":

```sh
cheat ls -t networking
```

To list all cheatsheets on the "personal" path:

```sh
cheat ls -p personal
```

To search for the phrase "ssh" among cheatsheets:

```sh
cheat search ssh
```

To search (by regex) for cheatsheets that contain an IP address:

```sh
cheat search -r '(?:[0-9]{1,3}\.){3}[0-9]{1,3}'
```

Flags may be combined in intuitive ways. Example: to search sheets on the
"personal" cheatpath that are tagged with "networking" and match a regex:

```sh
cheat search '(?:[0-9]{1,3}\.){3}[0-9]{1,3}' -p personal -t networking --regex
```

## Installing

```sh
# If you have mage installed:
mage

# Else:
go run mage.go
```

## Cheatsheets

Cheatsheets are plain-text files with no file extension, and are named
according to the command used to view them:

```sh
cheat view tar     # file is named "tar"
cheat view foo/bar # file is named "bar", in a "foo" subdirectory
```

Cheatsheet text may optionally be preceeded by a YAML frontmatter header that
assigns tags and specifies syntax:

```
---
syntax: javascript
tags: [ array, map ]
---
// To map over an array:
const squares = [1, 2, 3, 4].map(x => x * x);
```

The `cheat` executable includes no cheatsheets, but [community-sourced
cheatsheets are available][cheatsheets]. You will be asked if you would like to
install the community-sourced cheatsheets the first time you run `cheat`.

## Cheatpaths

Cheatsheets are stored on "cheatpaths", which are directories that contain
cheatsheets. Cheatpaths are specified in the `conf.yml` file.

It can be useful to configure `cheat` against multiple cheatpaths. A common
pattern is to store cheatsheets from multiple repositories on individual
cheatpaths:

```yaml
# conf.yml:
# ...
cheatpaths:
  - name: community                   # a name for the cheatpath
    path: ~/documents/cheat/community # the path's location on the filesystem
    tags: [ community ]               # these tags will be applied to all sheets on the path
    readonly: true                    # if true, `cheat` will not create new cheatsheets here

  - name: personal
    path: ~/documents/cheat/personal  # this is a separate directory and repository than above
    tags: [ personal ]
    readonly: false                   # new sheets may be written here
# ...
```

The `readonly` option instructs `cheat` not to edit (or create) any cheatsheets
on the path. This is useful to prevent merge-conflicts from arising on upstream
cheatsheet repositories.

If a user attempts to edit a cheatsheet on a read-only cheatpath, `cheat` will
transparently copy that sheet to a writeable directory before opening it for
editing.

### Directory-scoped Cheatpaths

At times, it can be useful to closely associate cheatsheets with a directory on
your filesystem. `cheat` facilitates this by searching for a `.cheat` folder in
the current working directory. If found, the `.cheat` directory will
(temporarily) be added to the cheatpaths.

## Autocompletion

```sh
cheat help completion
```

