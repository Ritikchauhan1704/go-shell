# go-shell

A simple shell implementation in Go.

## Getting Started

You can run this shell application using Docker. This ensures you have all the necessary dependencies and a consistent environment (Ubuntu).

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed on your machine.

### Build and Run

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/ShwetaRoy17/go-shell.git
    cd go-shell
    ```

2.  **Build the Docker image:**

    ```bash
    docker build -t go-shell .
    ```

3.  **Run the container:**

    To use the shell interactively, you must run the container with the `-it` flags:

    ```bash
    docker run -it --rm go-shell
    ```

    - `-i`: Keeps STDIN open so you can type commands.
    - `-t`: Allocates a pseudo-TTY.
    - `--rm`: Automatically removes the container when you exit.

<!-- "github.com/ShwetaRoy17/go-shell/app/shell" -->