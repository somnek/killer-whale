<p align="center">
  <img style="width:300px" src="https://github.com/somnek/killer-whale/blob/main/src/logo.png?raw=true"/>
</p>


# killer-whale 🐳

Killer Whale is a Docker TUI for terminal dwellers. It provides an intuitive, easy-to-use interface for managing your Docker containers without leaving the comfort of your command line.



## Usage

1. Clone the repository using Git: 

```bash
git clone https://github.com/somnek/killer-whale.git
```

2. Run the application using Go:

```bash
go run .
```
or build it:
```bash
cd killer-whale && go build -o killer-whale
```
Once the build is complete, move the executable to a directory in your system's `PATH` environment variable so that you can run it from anywhere.

For example, on Linux or macOS, you can move the executable to the `/usr/local/bin` directory:

3. Restart your terminal and run the application:

zsh:
```bash
source ~/.zshrc
```
bash:
```bash
exec bash
```

4. Run the application:

```bash
killer-whale
```

# Features
- [x] List all containers
- [x] Start/Stop containers
- [x] Restart containers
- [x] Remove containers
- [x] List all images

# WIP 🛠️
- [ ] docker logs -> full screen
- [ ] Remove images
- [ ] Hotkeys configuration

So why settle for a boring GUI when you can have a killer 🤘 command-line interface to manage your Docker containers? Killer Whale is designed to be fun and easy to use, with intuitive keyboard shortcuts ⌨️ and an attractive, streamlined interface that won't slow you down. With Killer Whale, you'll feel like a Docker pro in no time 🚀. Give it a try and see how killer your container management skills can be! 😎
