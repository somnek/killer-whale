<p align="center">
  <img style="width:300px" src="https://github.com/somnek/killer-whale/blob/main/src/logo.png?raw=true"/>
</p>

# killer-whale 🐳

Killer Whale is a Docker TUI for terminal dwellers. It provides an intuitive, easy-to-use interface for managing your Docker containers without leaving the comfort of your command line.

- more feature coming soon, thought its temping to add more feature, killer-whale meant to be as minimalistic as possible, if you are looking for docker tui with more features, check out [lazydocker](https://github.com/jesseduffield/lazydocker) by jesseduffield.

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
