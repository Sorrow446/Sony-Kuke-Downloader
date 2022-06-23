# Sony-Kuke-Downloader
[Sony Kuke](https://dereferer.me/?https://hi-resmusic.sonyselect.kuke.com/) downloader written in Go. This one was a pain, enjoy.
![](https://i.imgur.com/jF6FXQu.png)
[Windows, Linux, macOS and Android binaries](https://github.com/Sorrow446/Sony-Kuke-Downloader/releases)

# Setup
Active subscription required. Input Sony Select user ID into config file.
### Sony Select ID

~~1. Login on browser.~~

~~2. Right click, view page source.~~

~~3. Ctrl+f, search for `sonyselectid`.~~

It's now stored as a cookie.
1. Login on browser.
2. Right click, inspect.
3. Application -> Cookie, `USEIDENTIFY`.
![](https://i.imgur.com/TA9AoYK.png)

Configure any other options if needed.
|Option|Info|
| --- | --- |
|sonySelectId|Sony Select user ID.
|outPath|Where to download to. Path will be made if it doesn't already exist.
|trackTemplate|Track filename naming template. Vars: album, albumArtist, artist, title, track, trackPad, trackTotal, year.
|omitArtists|Omit album artists from album folder names.
|keepCover|Don't delete covers from album folders.

**FFmpeg is needed to concat FLAC segments.**    
[Windows (gpl)](https://github.com/BtbN/FFmpeg-Builds/releases)    
Linux: `sudo apt install ffmpeg`    
Termux `pkg install ffmpeg`

# Usage
Args take priority over the same config file options.

Download two albums:   
`sk_dl_x64.exe https://hi-resmusic.sonyselect.kuke.com/page/album.html?id=10628 https://hi-resmusic.sonyselect.kuke.com/page/album.html?id=10896`

Download a single album and from two text files:   
`sk_dl_x64.exe https://hi-resmusic.sonyselect.kuke.com/page/album.html?id=10628 G:\1.txt G:\2.txt`

```
 _____                _____     _          ____                _           _
|   __|___ ___ _ _   |  |  |_ _| |_ ___   |    \ ___ _ _ _ ___| |___ ___ _| |___ ___
|__   | . |   | | |  |    -| | | '_| -_|  |  |  | . | | | |   | | . | .'| . | -_|  _|
|_____|___|_|_|_  |  |__|__|___|_,_|___|  |____/|___|_____|_|_|_|___|__,|___|___|_|
              |___|

Usage: sk_dl_x64.exe [--outpath OUTPATH] URLS [URLS ...]

Positional arguments:
  URLS

Options:
  --outpath OUTPATH, -o OUTPATH
                         Where to download to. Path will be made if it doesn't already exist.
  --help, -h             display this help and exit
  ```
  
  # Disclaimer
- I will not be responsible for how you use Sony Kuke Downloader.    
- Sony brand and name is the registered trademark of its respective owner.    
- Sony Kuke Downloader has no partnership, sponsorship or endorsement with Sony.
