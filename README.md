# Wallpeek: Your Terminal's New Viewport 🖼️

**Wallpeek** is a blazing-fast, cross-platform terminal application that transforms your command line into a dynamic wallpaper browser and setter. Dive into your image collection without leaving the comfort of your terminal, and set your favorite wallpapers with a single keystroke!

## ✨ Features

*   **Blazing Fast:** Written in Go for unparalleled performance.
*   **Cross-Platform:** Works seamlessly on macOS and Linux.
*   **Terminal-Native:** Browse and set wallpapers directly from your CLI.
*   **Image Preview:** See your wallpapers rendered beautifully in your terminal (Kitty/iTerm2 compatible).
*   **Intuitive Keybindings:** Navigate, randomize, and set wallpapers with ease.

## 🚀 Installation

### From Source

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/x0d7x/wallpeek.git
    cd wallpeek
    ```
2.  **Build the application:**
    ```bash
    make build
    ```

### Cross-Platform Builds

You can also build Wallpeek for other operating systems using the provided `Makefile`:

*   **macOS:**
    ```bash
    make mac
    ```
*   **Linux:**
    ```bash
    make linux
    ```
*   **Windows:**
    ```bash
    make windows
    ```
3.  **Move to your PATH (optional):**
    ```bash
    sudo mv wallpeek /usr/local/bin/
    ```

### Dependencies

*   **`wypaper` (Linux only, optional):** For enhanced wallpaper setting on Linux. Install it via your distribution's package manager or from source.

## 💡 Usage

Simply run `wallpeek` followed by the path to your image directory:

```bash
wallpeek /path/to/your/images
```

### Keybindings

| Key           | Action                               |
| :------------ | :----------------------------------- |
| `j` / `↓`     | Next image                           |
| `k` / `↑`     | Previous image                       |
| `r`           | Random image                         |
| `s` / `Enter` | Set current image as wallpaper       |
| `q` / `Esc`   | Quit Wallpeek                        |

## 🤔 Why Wallpeek?

Tired of graphical file browsers just to pick a wallpaper? Wallpeek brings the power and speed of the terminal to your wallpaper management. It's perfect for minimalists, CLI enthusiasts, and anyone who wants a more efficient way to refresh their desktop.

## 🤝 Contributing

We welcome contributions! Feel free to open issues, submit pull requests, or suggest new features.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

