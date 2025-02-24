package main

import (
	// "fmt"
	"fmt"
	"math/rand"
	"sort"

	n "flappy/utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	Gravity      = 0.5
	FlapStrength = -6
	PipeWidth    = 70
	PipeGap      = 150
	PipeSpeed    = 3
	PipeSpacing  = 300
	ScrollSpeed  = 1
	MaxSpeed     = 5.0
	NumBirds     = 10
)

type Bird struct {
	X        float32
	Y        float32
	Velocity float32
	Alive    bool
	Size     float32
	Fitness  float32
	Brain    *n.NeuralNetwork
}

type Pipe struct {
	X      float32
	GapY   float32
	Offset float32
}

var birds []*Bird
var pipes []Pipe
var score int
var generations = 0

func evolve(population []*Bird) []*Bird {
	// Sort by fitness (descending order)
	sort.Slice(population, func(i, j int) bool {
		return population[i].Fitness > population[j].Fitness
	})

	// Select the top-performing birds (at least 1)
	topSize := len(population) / 10
	if topSize < 1 {
		topSize = 1
	}
	top := population[:topSize]

	// Create the next generation
	newPop := make([]*Bird, 0, len(population)) // Preallocate slice
	for len(newPop) < len(population) {
		parentA := top[rand.Intn(len(top))]
		var parentB *Bird

		// Ensure different parents are selected
		for {
			parentB = top[rand.Intn(len(top))]
			if parentB != parentA || len(top) == 1 {
				break
			}
		}

		// Crossover and mutate new child
		child := crossover(parentA, parentB)
		mutate(child)
		newPop = append(newPop, child)
	}

	return newPop
}

// Crossover: Mix two parents to create a new bird
func crossover(parentA, parentB *Bird) *Bird {
	// fmt.Println("ParentA brain weights: ", parentA.Brain.Weights1)
	// fmt.Println("ParentB brain weights: ", parentB.Brain.Weights1)
	child := &Bird{Brain: n.CreateNetwork(4, 6, 1)}
	for i := range child.Brain.Weights1 {
		for j := range child.Brain.Weights1[i] {
			if rand.Float64() > 0.5 {
				child.Brain.Weights1[i][j] = parentA.Brain.Weights1[i][j]
			} else {
				child.Brain.Weights1[i][j] = parentB.Brain.Weights1[i][j]
			}
		}
	}
	return child
}

// Mutation: Randomly adjust some weights
func mutate(bird *Bird) {
	for i := range bird.Brain.Weights1 {
		for j := range bird.Brain.Weights1[i] {
			if rand.Float64() < 0.05 { // 5% mutation rate
				bird.Brain.Weights1[i][j] += rand.NormFloat64() * 0.1
			}
		}
	}
}

func drawGame() {
	// Draw the bird
	for i := range birds {
		if birds[i].Alive {
			rl.DrawCircle(int32(birds[i].X), int32(birds[i].Y), birds[i].Size, rl.Yellow)
		} else {
			rl.DrawCircle(int32(birds[i].X), int32(birds[i].Y), birds[i].Size, rl.Red)
		}
	}

	// Draw the pipes
	for _, pipe := range pipes {
		// Top Rectangle
		rl.DrawRectangle(int32(pipe.X), 0, PipeWidth, int32(pipe.GapY), rl.Green)

		// Bottom Rectangle
		rl.DrawRectangle(int32(pipe.X), int32(pipe.GapY+PipeGap), PipeWidth, ScreenHeight, rl.Green)
	}
}

// Check if the bird collides with a specific pipe
func isCollisionWithPipe(pipe *Pipe, bird *Bird) bool {
	// Check if the bird is within the horizontal bounds of the pipe
	if bird.X+bird.Size > pipe.X && bird.X-bird.Size < pipe.X+PipeWidth {
		// Check if the bird is either above the gap (top pipe) or below the gap (bottom pipe)
		if bird.Y-bird.Size < pipe.GapY || bird.Y+bird.Size > pipe.GapY+PipeGap {
			return true // Bird collides with the pipe
		}
	}
	return false
}

func updateGame() {
	numDead := 0
	for i := range birds {
		birds[i].Y += birds[i].Velocity
		if birds[i].Velocity < MaxSpeed {
			birds[i].Velocity += Gravity
		}

		if birds[i].Y >= ScreenHeight-30 {
			birds[i].Alive = false
		}

		// Increase fitness or kill the bird
		if birds[i].Alive {
			birds[i].Fitness += 1
		} else if !birds[i].Alive {
			numDead += 1
		}

		// Create our inputs to try and predict

		// Distance to current pipe
		// Distance to the pipe gap
		// Y position and velocity

		// pipe_dist := pipes[0].X
		// fmt.Println("Pipe distance: ", pipe_dist)
		//
		// 	birds[i].Brain.Predict()

	}

	for i := range pipes {
		pipes[i].X -= ScrollSpeed

		if pipes[i].X < -PipeWidth {
			// Create a new pipe, and add it to the end of the list
			newPipe := pipes[i]
			newPipe.X = ScreenWidth                                          // Reset pipe to right side
			newPipe.GapY = float32(rand.Intn(ScreenHeight-PipeGap-100) + 50) // Randomize gap position

			pipes = pipes[1:]
			pipes = append(pipes, newPipe)
		}

		// Check for collisions
		for j := range birds {
			if isCollisionWithPipe(&pipes[i], birds[j]) {
				birds[j].Alive = false
			}
		}
	}

	// All birds have died, so we want to reset the game!
	if numDead == NumBirds {

		generations += 1
		resetGame()
		// fmt.Println("Generations: ", generations)
	}
}

func resetGame() {
	fmt.Println("Resetting game!")

	// Initial game load

	// Otherwise, we want to crossover the birds!
	birds = evolve(birds)

	for i := range birds {
		*birds[i] = Bird{
			X:        30,
			Y:        ScreenHeight / 2,
			Velocity: 0,
			Alive:    true,
			Size:     20,
			Fitness:  0,
			Brain:    birds[i].Brain,
		}

	}

	distance := 200.0
	pipes = make([]Pipe, 3)
	for i := range pipes {

		pipes[i] = Pipe{
			X:      float32(distance),
			GapY:   float32(rand.Intn(ScreenHeight-PipeGap-100) + 50),
			Offset: float32(distance - 200.0),
		}

		distance += PipeSpacing
	}
}

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Flappy Bird")
	rl.SetTargetFPS(60)

	for i := 0; i < NumBirds; i++ {
		birds = append(birds, &Bird{
			X:        30,
			Y:        ScreenHeight / 2,
			Velocity: 0,
			Alive:    true,
			Size:     20,
			Brain:    n.CreateNetwork(4, 6, 1),
		})
	}

	distance := 200.0
	for i := 0; i < 3; i++ {

		pipes = append(pipes, Pipe{
			X:      float32(distance),
			GapY:   float32(rand.Intn(ScreenHeight-PipeGap-100) + 50),
			Offset: float32(distance - 200.0),
		})

		distance += PipeSpacing

	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		updateGame()

		drawGame()

		if rl.IsKeyDown(rl.KeySpace) {
			birds[0].Velocity = FlapStrength
		}

		rl.EndDrawing()
	}
}
