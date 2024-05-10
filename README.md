# Moorhuhn Mods

An all-in-one experimental Go command-line tool for patching Moorhuhn data files.

## Features

- [x] Unpacking data files;
- [x] Repacking data files;
- [x] Packing and merging multiple mods into unified data file;
- [x] A Python webserver, to patch and download modded data files on the fly;

## Supported games

- [x] Moorhuhn Kart Extra XXL (game ID: `mhk_1`/`mhk_extra`) - untested!;
- [x] Moorhuhn Kart 2 (game ID: `mhk_2`);
  - [ ] Schatzjäger (untested) (game ID: `schatzjaeger`) - note: this uses the exact same algorithm as MHK2;
- [x] Moorhuhn Kart 3 (game ID: `mhk_3`);
- [x] Moorhuhn Kart 4 (Thunder) (game ID: `mhk_4`/`mhk_thunder`);

## Wait, what?

Moorhuhn Kart games are the games of my childhood. I've always wanted to mod them, however it seems that there are a very few enthusiasts like me who are interested in modding these games. In fact, I thought that I was the only one.

Later, I've randomly stumbled upon the [QuickBMS](https://aluigi.altervista.org/quickbms.htm) scripts for unpacking and repacking data files. I've quickly bashed an (now old) [repository](https://github.com/SKevo18/mhk_mods) together, which was a simple, but clunky Python CLI tool for unpacking and repacking data files. It also included a webserver for building, merging and downloading modded data files on the fly for casual people that are not familiar with programming or CLI tools.

Then, I have discovered the [Moorhuhnverse Discord server](https://discord.gg/buJ64SrHxY) which is a community of fellow Moorhuhn enthusiasts. It turn out that there are some very talented people who modded a few Moorhuhn games.

By exchanging knowledge and helpful tips, I've been able to understand how the data files are structured and how to modify them (before, I had no idea how to do that, and would only plug existing tools together). As the repository grew, stacked on top of older and older code, I've decided to rewrite the whole thing in pure Go, to make it more maintainable, without the need to rely on QuickBMS scripts and Python.

## Usage

### Download existing mods

The [webserver](./webserver/) is a simple Python webserver, which can be used to build and download modded data files on the fly, without the need for end-users to use the command-line tool at all. See it live at [https://mhmods.svit.ac/](https://mhmods.svit.ac/).

### Download

You can obtain the latest version of the tool for your operating system from the [releases](https://github.com/SKevo18/mh_mods/releases) page. They are cross-compiled for Windows, Linux and MacOS (GitHub Actions).

### Compile

If you want to compile the tool yourself, you must have Go 1.22.0 installed on your system. Then, you can run the following command to compile the tool:

```bash
# clone the repository and `cd` to repo root, then:
go build -o build/
```

The binary for your OS will be available in the `build/` directory. You can also compile the tool and run it right away by using `go run .` in the repository root.

## FaQ

### What data files are you talking about?

Moorhuhn games use a special data file format, which contains all the game assets, such as textures, models, sounds, configuration files, etc. The data files are present in your Moorhuhn game directory, and have the `.dat` (MHK 1, 2) or `.sar` (MHK 3, 4) extension.

These data files are packed with a custom compression algorithm, and are not directly readable by any standard tools. That's why a special tool is needed to unpack and repack them, which is what this repository is about.

### What can I do with the tool?

While this tool is not a full-fledged modding suite, it provides the basic functionality to unpack and repack data files, and to merge multiple mods into a single data file.

You can't use the tool to create new assets, such as models or textures, but you can modify existing ones. For example, you can replace the textures of the game characters, modify existing levels/tracks or modify the configuration files to change the game behavior.

The tools aims to simplify this process as much as possible, and to provide a simple way to share mods with others. It also includes a webserver, which can be used to build and download modded data files on the fly, without the need for end-users to use the command-line tool at all.

## Credits

This tool wouldn't be possible without the help of a few talented humans. I'd like to thank the following people for their help towards making this tool a reality:

- **Luigi Auriemma**, for creating the original QuickBMS tool that sparked my interest in (Moorhuhn) modding;
- **pyramidensurfer**, for their incredible dedication towards understanding the Moorhuhn data file format, and for providing example code and tools to work with the data files;
- **Blue Cap guy**, for their valuable insights and feedback on the tool;
- The entire **Moorhuhnverse Discord community** and other fellow Moorhuhn enthusiasts, for preserving the interest in Moorhuhn games;

Thank you all very much!
