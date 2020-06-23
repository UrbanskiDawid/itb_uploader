# ITB UPLOADER
[![](https://github.com/UrbanskiDawid/itb_uploader/workflows/build/badge.svg)](https://github.com/UrbanskiDawid/itb_uploader/actions?query=workflow%3Abuild)[![Go Report Card](https://goreportcard.com/badge/github.com/UrbanskiDawid/itb_uploader)](https://goreportcard.com/report/github.com/UrbanskiDawid/itb_uploader)

### REST api for ssh commands & files

# WARNING

WIP, do not use in production!
code quality: bad, this is my first experiment with GO lang


## how to use

Start by editing **config.json** to fit your needs.
You have to: define all [remote servers](#Configuration-servers) and define all [actions](#Configuration-actions)

then [start-as-server](#start-as-server) on remote machine/machines and use [command-line](#start-as-client) to execute commands.

### Configuration

There are two sections **servers** and **actions**
Firts defines where actions are executed second what actions shall be done. 

```
{
    "servers": [],
    "actions": []
}
```

#### Configuration - servers
```
    "servers": [
    {
        "nickname": "mini",
        "host": "192.168.1.111",
        "auth":{ "user": "userA", "pass": "passA"},
        "port": 22
    },
```
should be self explanatory.

NOTE: you can use path to ssh key as "pass"

#### Configuration - actions

there are 3 types of actions: 
- command
- file upload
- file download

##### ACTIONS -  command
in **config.json** in *actions* define:

```
    {
        "name": "voice",
        "cmd": "python3 /home/dawid/cast.py ASSISTANT_VOICE/dave.mp3",
        "server": "MINI"
    },
```

##### ACTIONS -  file upload
in **config.json** in *actions* define:

```
    {
        "name": "upload",
        "fileTarget": "/home/dave/test",
        "server": "localhost"
    },
```

##### ACTIONS -  file download
in **config.json** in *actions* define:

```
    {
        "name": "download",
        "fileDownload": "/home/dave/GIT/GO/go.mod",
        "server": "localhost"
    },
```

## start as client
```$ ./itb_uploader.linux help```


## start as server
```$ ./itb_uploader.linux server```
