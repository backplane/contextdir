# contextdir

Utility for working with local Docker build context directories


## Development notes

This program is a work in progress. Some of the planned features are implemented but output formats (particularly the `--detailed` ones) are very likely to change.

Planned features:

1. ✅ list the context dir (. by default)
2. ✅ give the user the ability to list only ignored files (show which files exist and are ignored, don't just print the ignore file)
3. ✅ calculate and print the checksum of the whole build context
4. compare the current context dir to the tags on a public image
5. maybe add a mode where it can print the top 3 directories and files by size to make it easier to find problems

## Usage

The following text is produced when the program is invoked with the `-h` or `--help` argument:

```
NAME:
   contextdir - Utility for working with local Docker build context directories

USAGE:
   contextdir [global options] command [command options] [arguments...]

VERSION:
   dev

COMMANDS:
   list, ls             list the contents of the given context dir, honoring .dockerignore (if found)
   checksum, sum, hash  list the contents of the given context path(s), excluding any paths mentioned in
                        .dockerignore
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbosity value  Sets the verbosity level of the log messages printed by the program, should be
      one of: "debug", "error", "fatal", "info", "panic", "trace", or "warn"
   --help, -h     show help
   --version, -v  print the version
```

### `list` subcommand

The following text is produced when invoking the `list` subcommand with the `-h` or `--help` arguments:

```
NAME:
   contextdir list - list the contents of the given context dir, honoring .dockerignore (if found)

USAGE:
   contextdir list [command options] [dir]

OPTIONS:
   --ignored  Instead of the normal output, only list files in the given contextdir(s) which
      match patterns in the .dockerignore file (default: false)
   --detailed  Format the output in JSON format with detailed information about the files (may
      be combined with the --ignored flag) (default: false)
   --help, -h  show help
```

If no `dir` argument is given the program defaults to the current working directory (`.`).

When working in a docker build context directory the list command prints a list of files that are in the build context that aren't excluded by the `.dockerignore` file:

```sh
$ contextdir ls
.dockerignore
Dockerfile
chromium_ssb.sh
entrypoint.sh
fonts.conf
```

The files that *are* ignored can be listed by using the `--ignored` flag:

```sh
$ contextdir ls --ignored
README.md
seccomp.json
```

Detailed information about the listed files can be obtained by adding the `--detailed` flag:

```sh
$contextdir ls --ignored --detailed
[
 {
  "Stats": {
   "path": "README.md",
   "mode": 420,
   "uid": 501,
   "gid": 20,
   "size": 4857,
   "modTime": 1677011647176224836
  },
  "StatsHash": "82f5ab9a5d24b61f02f281723e4f19868b682f551bb2c7e898424eb366d16bcc"
 },
 {
  "Stats": {
   "path": "seccomp.json",
   "mode": 420,
   "uid": 501,
   "gid": 20,
   "size": 36373,
   "modTime": 1637673846161641661
  },
  "StatsHash": "337375d5b7873261b615f82af96d9519b0e5cd5f009ff493027f8c5baf130ce4"
 }
]
```

### `checksum` subcommand

The following output is produced when invoking the checksum subcommand with the `-h` or `--help` arguments:

```
NAME:
   contextdir checksum - list the contents of the given context path(s), excluding any paths mentioned in
                         .dockerignore

USAGE:
   contextdir checksum [command options] [dir]

OPTIONS:
   --json      Format the output in JSON format (default: false)
   --detailed  Report the entire scan in JSON format (default: false)
   --help, -h  show help
```

If no `dir` argument is given the program defaults to the current working directory (`.`).

The checksum of the file statistics of all non-ignored files in the build context directory is printed when the `checksum` command is used:

```sh
$ contextdir checksum
662b056325942f309289f18efad7c13098cce2d7803255a67c07816124538579
```

This checksum should be stable until any of the included stats of the included files change. The standard report is availabe in JSON format with the `--json` flag:

```sh
$ contextdir checksum --json
{
 "Path": ".",
 "Checksum": "662b056325942f309289f18efad7c13098cce2d7803255a67c07816124538579"
}
```

A very detailed reporting of the information used to calculate the checksum is available via the `--detailed` flag:

```sh
$ contextdir checksum --detailed
{
  "Dir": {
    "path": ".",
    "files": [
      {
        "Stats": {
          "path": ".dockerignore",
          "mode": 420,
          "uid": 501,
          "gid": 20,
          "size": 23,
          "modTime": 1677002480364561200
        },
        "StatsHash": "4f9bbd8fcca5a2fd1da4cc543cc080979130eebb27775948f7247bda088709e0"
      },
      {
        "Stats": {
          "path": "Dockerfile",
          "mode": 420,
          "uid": 501,
          "gid": 20,
          "size": 1186,
          "modTime": 1677007954555730700
        },
        "StatsHash": "1f6ed07c16794923ff44f59717a863c96732ddf67b26593aa33898b9bfde2b9b"
      },
      {
        "Stats": {
          "path": "chromium_ssb.sh",
          "mode": 493,
          "uid": 501,
          "gid": 20,
          "size": 5916,
          "modTime": 1677024828470154000,
          "xattrs": {
            "com.apple.TextEncoding": "dXRmLTg7MTM0MjE3OTg0"
          }
        },
        "StatsHash": "c9fe14fef5f7078ae6327ba11949e3cfee8a4a4f707f15576839c28acf1c21ff"
      },
      {
        "Stats": {
          "path": "entrypoint.sh",
          "mode": 493,
          "uid": 501,
          "gid": 20,
          "size": 925,
          "modTime": 1677002608063513300
        },
        "StatsHash": "1ca9e14b99f28136111262ce8cfc46a3f46533733d7626554ec020dfe730bd28"
      },
      {
        "Stats": {
          "path": "fonts.conf",
          "mode": 420,
          "uid": 501,
          "gid": 20,
          "size": 676,
          "modTime": 1637673846161205000
        },
        "StatsHash": "59692f2172a20cb4b850b4acb2ece5002c0d4c90536e75f18439422edb3ce018"
      }
    ],
    "hash": "662b056325942f309289f18efad7c13098cce2d7803255a67c07816124538579"
  },
  "DockerIgnoreEntries": [
    "README.md",
    "seccomp.json"
  ],
  "IgnoredFiles": [
    {
      "Stats": {
        "path": "README.md",
        "mode": 420,
        "uid": 501,
        "gid": 20,
        "size": 4857,
        "modTime": 1677011647176224800
      },
      "StatsHash": "82f5ab9a5d24b61f02f281723e4f19868b682f551bb2c7e898424eb366d16bcc"
    },
    {
      "Stats": {
        "path": "seccomp.json",
        "mode": 420,
        "uid": 501,
        "gid": 20,
        "size": 36373,
        "modTime": 1637673846161641700
      },
      "StatsHash": "337375d5b7873261b615f82af96d9519b0e5cd5f009ff493027f8c5baf130ce4"
    }
  ],
  "FoundDockerfile": true,
  "FoundDockerignore": true
}
```


The key .Dir.hash is the SHA-256 sum of all the StatsHash values of every non-ignored file, this is the value that is normally reported without the `--detailed` flag.

