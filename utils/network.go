package utils

import (
	"math"
	"math/rand"
	// "math/rand"
)

// Sigmoid function squashes large values to a number between 0 and 1
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

type NeuralNetwork struct {
	InputSize, HiddenSize, OutputSize int
	Weights1, Weights2                [][]float64
}

func randomMatrix(rows, cols int) [][]float64 {
	matrix := make([][]float64, rows) // Create all the rows
	for i := range matrix {
		matrix[i] = make([]float64, cols) // Create the column

		for j := range matrix[i] {
			// Random values between -1 and 1
			matrix[i][j] = rand.Float64()*2 - 1
		}
	}

	return matrix
}

func CreateNetwork(input, hidden, output int) *NeuralNetwork {
	nn := &NeuralNetwork{
		InputSize:  input,
		HiddenSize: hidden,
		OutputSize: output,
		Weights1:   randomMatrix(input, hidden),
		Weights2:   randomMatrix(hidden, output),
	}

	return nn
}

// Take the matrix, and do a forward pass through the weights
func (nn *NeuralNetwork) Predict(inputs []float64) float64 {
	// Go through the hidden layer and sum the weights
	hidden := make([]float64, nn.HiddenSize)

	for i := 0; i < nn.HiddenSize; i++ {
		sum := 0.0
		for j := 0; j < nn.InputSize; j++ {
			sum += inputs[j] * nn.Weights1[j][i]
		}

		hidden[i] = sigmoid(sum)
	}

	// Computer output layer activation
	output := 0.0
	for i := 0; i < nn.HiddenSize; i++ {
		output += hidden[i] * nn.Weights2[i][0]
	}

	// Probability of flapping
	return sigmoid(output)
}
