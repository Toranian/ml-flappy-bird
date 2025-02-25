package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	n "flappy/utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	GameSpeed    = 1.5
	ScreenWidth  = 800
	ScreenHeight = 600
	Gravity      = 0.5 * GameSpeed
	FlapStrength = -6 * GameSpeed
	PipeWidth    = 70
	PipeGap      = 200
	PipeSpeed    = 3 * GameSpeed
	PipeSpacing  = 300
	ScrollSpeed  = 1 * GameSpeed
	MaxSpeed     = 5.0 * GameSpeed
	NumBirds     = 30
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
var score = 0
var generations = 0
var pipesPassed = 0
var start = time.Now()
var highScore = 0

func evolve(population []*Bird) []*Bird {
	// Sort by fitness (descending order)
	sort.Slice(population, func(i, j int) bool {
		return population[i].Fitness > population[j].Fitness
	})

	// Select the top-performing birds (at least 1)
	topSize := int(math.Floor(float64(len(population)) * 0.3))
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
	child := &Bird{Brain: n.CreateNetwork(5, 8, 1)}
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
	mutationRate := 0.1 / float64(generations+1)
	for i := range bird.Brain.Weights1 {
		for j := range bird.Brain.Weights1[i] {
			if rand.Float64() < mutationRate { // 5% mutation rate
				bird.Brain.Weights1[i][j] += rand.NormFloat64() * 0.1
			}
		}
	}
}

func drawGame() {
	// Draw the bird
	for i := range birds {
		if birds[i].Alive {
			rl.DrawCircle(int32(birds[i].X), int32(birds[i].Y), birds[i].Size, rl.Orange)
		}
	}

	// Draw the pipes
	for _, pipe := range pipes {
		// Top Rectangle
		rl.DrawRectangle(int32(pipe.X), 0, PipeWidth, int32(pipe.GapY), rl.Green)

		// Bottom Rectangle
		rl.DrawRectangle(int32(pipe.X), int32(pipe.GapY+PipeGap), PipeWidth, ScreenHeight, rl.Green)
	}

	// Stats
	// for count := range birds {
	// 	rl.DrawText(fmt.Sprintf("Bird %d fitness: %f", count+1, birds[count].Fitness), 20, int32(20*count+1), 15, rl.Black)
	// }

	rl.DrawText(fmt.Sprintf("%d", score), ScreenWidth-60, 20, 30, rl.Black)
	rl.DrawText(fmt.Sprintf("%d", highScore), ScreenWidth-60, 60, 20, rl.Black)
	diff := time.Now().Sub(start)
	rl.DrawText(fmt.Sprintf("Generation: %d \nTime: %.2fs", generations, diff.Seconds()), 20, 20, 20, rl.Black)

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

		if birds[i].Alive && (birds[i].Y >= ScreenHeight+15 || birds[i].Y <= 0) {
			birds[i].Alive = false
			birds[i].Fitness -= 50
		}

		// Increase fitness or kill the bird
		if birds[i].Alive {
			birds[i].Fitness += 0.5
		} else if !birds[i].Alive {
			numDead += 1
			continue
		}

		// Create our inputs to try and predict
		pipe_dist := pipes[0].X
		gap_dist_top := pipes[0].GapY - birds[i].Y
		gap_dist_bottom := pipes[0].GapY + PipeGap - birds[i].Y

		// If between the pipes, we want to increase the fitness
		if birds[i].X-29 > pipes[0].X+PipeWidth {
			pipesPassed += 1
			birds[i].Fitness += 100

			if pipesPassed > score {
				score = pipesPassed
			}
		}

		// rl.DrawLine(30, int32(birds[i].Y), int32(pipes[0].X), int32(pipes[0].GapY), rl.Purple)
		// rl.DrawLine(30, int32(birds[i].Y), int32(pipes[0].X), int32(pipes[0].GapY+PipeGap), rl.Red)

		input := []float64{
			float64(birds[i].Velocity),
			float64(pipe_dist),
			float64(gap_dist_top),
			float64(gap_dist_bottom),
			float64(birds[i].Y),
		}
		activate := birds[i].Brain.Predict(input)

		if activate >= 0.5 {
			birds[i].Velocity = FlapStrength
		}
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
			if isCollisionWithPipe(&pipes[i], birds[j]) && birds[j].Alive {
				birds[j].Alive = false

				// Reward birds for being closer to the gap when they die
				distanceToPipe := ((ScreenHeight/2)/(pipes[i].GapY+(PipeGap/2)) - birds[j].Y) / 10
				birds[j].Fitness -= distanceToPipe
			}
		}
	}

	// All birds have died, so we want to reset the game!
	if numDead == NumBirds {

		generations += 1
		resetGame()
	}
}

func resetGame() {

	if score > highScore {
		highScore = score
	}
	score = 0

	diff := time.Now().Sub(start)
	fmt.Printf("Generation: %d (%.2fs)\n", generations, diff.Seconds())
	start = time.Now()

	// Otherwise, we want to crossover the birds!
	birds = evolve(birds)

	for i := range birds {
		*birds[i] = Bird{
			X:        30,
			Y:        ScreenHeight / 2,
			Velocity: 0,
			Alive:    true,
			Size:     20,
			Fitness:  birds[i].Fitness,
			Brain:    birds[i].Brain,
		}

	}

	distance := 200.0
	pipes = make([]Pipe, 3)
	for i := range pipes {

		pipes[i] = Pipe{
			X:      float32(distance + 300),
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
			Brain:    n.CreateNetwork(5, 8, 1),
		})
	}

	distance := 200.0
	for i := 0; i < 3; i++ {

		pipes = append(pipes, Pipe{
			X:      float32(distance + 300),
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

	highestScore := 0
	bestBird := birds[0]
	for i := range birds {
		if birds[i].Fitness > float32(highestScore) {
			highestScore = int(birds[i].Fitness)
			bestBird = birds[i]
		}
	}

	fmt.Println("Best fitness: ", highestScore)
	fmt.Println("Highest Pipes Passed", score)
	fmt.Println("\nBird weights: ", bestBird.Brain)

}
