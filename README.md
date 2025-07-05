# note

`note` allows you to create and view interactive notes on the command-line. It
was designed to help remind \*nix system administrators of options for commands
that they use frequently, but not frequently enough to remember.

![The obligatory xkcd](http://imgs.xkcd.com/comics/tar.png 'The obligatory xkcd')

## Example

The next time you're forced to disarm a nuclear weapon without consulting
Google, you may run:

```sh
note v tar
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
tags: ["note"]
---

# View a note:
note view tar
note v tar

# View a nested note called bar:
note view foo/bar


# Opens the "tar" notesheet for editing, or creates it if it does not exist:
note edit tar

# Nested notesheets are accessed like this:
note edit foo/bar
note e foo/bar


# To list the notebooks:
note books
note b


# To list all available notes:
note ls

# To list all notesheets that are tagged with "networking":
note ls -t networking

# To list all notesheets on the "personal" path:
note ls -p personal


# To search for the phrase "ssh" among notesheets:
note search ssh
note s ssh

# To search (by regex) for notesheets that contain an IP address:
note search -r '(?:[0-9]{1,3}\.){3}[0-9]{1,3}'

# Flags may be combined:
note search '(?:[0-9]{1,3}\.){3}[0-9]{1,3}' -p personal -t networking --regex
```

Yes, this is also a note.

## Installing

```sh
gh repo clone yagoyudi/note
go build .
```

Or install a binary from the releases.

## Notebook

Notes are stored on notebooks, which are directories that contain notes.
Notebooks are specified in the `config.yml` file.

It can be useful to configure `note` against multiple notebooks. A common
pattern is to store notes from multiple repositories on individual notebooks:

```yaml
# config.yml:
# ...
notebooks:
  - name: community
    # The path's location on the filesystem:
    path: ~/documents/note/community
    # These tags will be applied to all sheets on the path:
    tags: [ community ]
    # If true, `note` will not create new note here:
    readonly: true

  - name: personal
    path: ~/documents/note/personal
    tags: [ personal ]
    readonly: false
# ...
```

The `readonly` option instructs `note` not to edit (or create) any note on the
notebook. This is useful to prevent merge-conflicts from arising on upstream
note repositories.

If a user attempts to edit a note on a read-only notebook, `note` will
transparently copy that note to a writeable directory before opening it for
editing.

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

You will be asked if you would like to install the community-sourced notebooks
the first time you run `note`.

## Autocompletion

```sh
note help completion
```
