# EmuDeck ROM Uploader

Upload ROMs for emulation to the SteamDeck easily over the air using sftp.

> [!Note]
> SSH needs to be enabled on the Steamdeck. A tutorial on how to do that can be found [here](https://www.youtube.com/watch?v=IWgJrrrQn6I)

## Install using Go

```sh
go install github.com/Lofter1/edru@latest
```

## Basic usage

```sh
edru --system psx --rom "/my/rom/path"
```

If the path to the ROM is a directory, the directory will be uploaded recursively.
If the path is a single file, only that file will be uploaded.

The files will be uploaded to the EmuDeck ROM folder for the specified system.
