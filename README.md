# Cheat

`cheat` allows you to create and view interactive notes on the command-line. It
was designed to help remind \*nix system administrators of options for commands
that they use frequently, but not frequently enough to remember.

![The obligatory xkcd](http://imgs.xkcd.com/comics/tar.png 'The obligatory xkcd')

Use `cheat` with [cheatsheets](https://github.com/cheat/cheatsheets).

## Example

The next time you're forced to disarm a nuclear weapon without consulting
Google, you may run:

```sh
cheat view tar
```

You will be presented with a note resembling the following:

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

```sh
---
syntax: bash
tags: ["cheat"]
---

# View a note:
cheat view tar

# View a nested note called bar:
cheat view foo/bar

# Alias for view:
cheat v tar


# Opens the "tar" cheatsheet for editing, or creates it if it does not exist:
cheat edit tar

# Nested cheatsheets are accessed like this:
cheat edit foo/bar

# Alias for edit:
cheat e foo/bar


# To view the configured cheatpaths:
cheat dirs


# To list all available cheatsheets:
cheat ls

# To list all cheatsheets that are tagged with "networking":
cheat ls -t networking

# To list all cheatsheets on the "personal" path:
cheat ls -p personal


# To search for the phrase "ssh" among cheatsheets:
cheat search ssh

# To search (by regex) for cheatsheets that contain an IP address:
cheat search -r '(?:[0-9]{1,3}\.){3}[0-9]{1,3}'

# Flags may be combined:
cheat search '(?:[0-9]{1,3}\.){3}[0-9]{1,3}' -p personal -t networking --regex

# Alias for search:
cheat s ssh
```

This is also a note.

## Installing

```sh
gh repo clone yagoyudi/cheat
go build .
```

Or install a binary from the releases.

## Notes

Notes are plain-text files with no file extension.

Notes may optionally be preceeded by a YAML header that assigns tags and
specifies syntax:

```note
---
syntax: javascript
tags: [ array, map ]
---

// To map over an array:
const squares = [1, 2, 3, 4].map(x => x * x);
```

The `cheat` executable includes no notebooks, but [community-sourced
cheatsheets are available](https://github.com/cheat/cheatsheets). You will be
asked if you would like to install the community-sourced cheatsheets the first
time you run `cheat`.

## Notebook

Notes are stored on "notebooks", which are directories that contain notes.
Notebooks are specified in the `conf.yml` file.

It can be useful to configure `cheat` against multiple notebooks. A common
pattern is to store notes from multiple repositories on individual notebooks:

```yaml
# conf.yml:
# ...
notebooks:
  - name: community
    # The path's location on the filesystem:
    path: ~/documents/cheat/community
    # These tags will be applied to all sheets on the path:
    tags: [ community ]
    # If true, `cheat` will not create new note here:
    readonly: true

  - name: personal
    path: ~/documents/cheat/personal
    tags: [ personal ]
    readonly: false
# ...
```

The `readonly` option instructs `cheat` not to edit (or create) any note on the
notebook. This is useful to prevent merge-conflicts from arising on upstream
note repositories.

If a user attempts to edit a note on a read-only notebook, `cheat` will
transparently copy that note to a writeable directory before opening it for
editing.

## Autocompletion

```sh
cheat help completion
```
