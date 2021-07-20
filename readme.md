# libman

Libman is an interactive spotify shell.

# Features

-	Control your spotify playback.
-	Edit playlists.
-	Fully complies to the spotify web api terms of usage.
-	Set aliases for commands.
-	Execute commands from a file during startup.
-	No need for premium.
-	Simple to get started, comfortable to use.
-	Has context aware in-application tab completions!

## Authors Note

Libman is still in development stages, the usage won't change but a lot of new features are on their way, so if you can't do `X` yet, just wait for it.

# Installation

Libman is written in `go`, and will work only with go1.16 and above so get an up to date go compiler.

You can either clone the repo or install directly with the go command:

```sh
# recommended:
git clone https://github.com/insomnimus/libman
cd libman
git checkout main
go install

# you can also do this:
# go install github.com/insomnimus/libman@latest
```

# Getting Started

You don't need a premium account, but you need to get an app token; visit [this link](https://developer.spotify.com/documentation/web-api/).

After registering an application, note down your client id, secret and redirect uri (you configure your own redirect uri).

The redirect URI should be a loopback (localhost); for example:

`http://localhost:8080/callback`

Now you can either launch libman to generate a config file and save your credentials there, or you can 
set some env variables:

-	`LIBMAN_ID`: Set it to your spotify client ID.
-	`LIBMAN_SECRET`: Set it to your spotify client secret.
-	`LIBMAN_REDIRECT_URI`: Set it to the redirect URI you configured from the spotify developer portal.
-	`LIBMAN_CACHE_PATH`:  This is not necessary but helpful, set it to a file where your session token will be saved so you won't have to authorize each time you launch libman.

That's it! Just launch libman and enjoy some music.

# Tips & Tricks

-	Create a `~/.libmanrc` file and write any valid libman command to be ran at startup (define your aliases here).
-	You can change the prompt! put `prompt MyNewPrompt>` in the libman config file.
-	The config file is located at `~/.config/libman.toml` on unix, and `C:\Users\username\AppData\Roaming` on windows.
