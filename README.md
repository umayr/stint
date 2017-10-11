# `stint`
> a super tiny worker that runs in background and fetches torrent information from RSS feeds

## Install

```
# if you have go setup in your system
λ go get -u github.com/umayr/stint/cmd/stint
  
# build it with source
λ git clone https://github.com/umayr/stint && cd $_ && make build
  
```

## Usage

```
λ stint --help                                                                                                                       (stint) 18:39:18
NAME:
   stint - a super tiny worker that runs in background and fetches torrent information from RSS feeds

USAGE:
   stint [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level value, -l value    set logging level, available options are debug, info, warn, error, and fatal
   --config-file value, -f value  path to configuration file
   --help, -h                     show help
   --version, -v                  print the version

```

Create a configuration file `.stintrc` in your home directory like:
```yaml
# RSS Feed URL
url: 'https://eztv.ag/ezrss.xml'
# command that would be executed once there's match
cmd: echo
# arguments that would be provided to the command
args: '{{ .Title }}'
# add filters below, for example:
shows:
  -
     title: 'Rick and Morty' # title of the show, it should be as clear as possible to avoid conflicts
     quality: high # it could be either normal, medium or high
```

And that's it, execute `stint` however you like, manually or set up a cron to check the RSS periodically. 
