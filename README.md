# Flappy Bird in Go + Raylib. AI from Scratch

[demo.webm](https://github.com/user-attachments/assets/10a98e98-b76a-402e-94aa-39add5bec3de)

## Approach

### Neural Network

Each bird in the program has it's own neural network. This neural network is as follows:

- 5 input nodes
- 1 hidden layer with 8 nodes
- 1 output node

#### Input

The 5 inputs are important to properly training the model, and are as follows:

1. The birds Y position
1. The position to the closest pipe (X distance)
1. The distance to the top of the pipe gap (Y distance)
1. The distance to the bottom of the pipe gap (Y distance)
1. The velocity

#### Output

If the output node is > 0.5, the bird flaps its wings.

### Fitness

## Install and Setup

1. Clone the code using the green button in the top right of the GitHub interface.
   1. `git clone git@github.com:Toranian/ml-flappy-bird.git`
1. In your terminal, install the Go libraries and Raylib:
   1. `go mod tidy`
   1. This project uses the Go bindings for Raylib, but that means you still have to install Raylib as well. To do so, you can follow the guide here: [gen2brain/raylib-go](https://github.com/gen2brain/raylib-go)
1. Once everything is installed, you can run the project with `go run ./main.go`.

! Note: Running for the first time will take a few minutes. This is because Raylib is being _compiled_ in order to run the game. After the first run through, it's very quick.
