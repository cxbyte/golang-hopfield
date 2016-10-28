package main

import (
	"os"
	"strconv"
	"reflect"
	"io/ioutil"
	"bufio"
)

type Matrix [][]int64;
type Vector []int64;
const MAX_CONVERGENCE_STEPS = 10;
const DIR_SAMPLES = "samples";
const DIR_SAMPLES_ORIGINAL = DIR_SAMPLES + "/original";

/**
Make vector
 */
func makeVector(sample string) (Vector) {
	file, _ := os.Open(sample);

	vector := Vector{};
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, v := range scanner.Text(){
			var vectorInt int64;
			if string(v) == "x" {
				vectorInt = 1;
			} else {
				vectorInt = -1;
			}

			vector = append(vector, vectorInt);
		};
	}

	return vector;
}

/**
Multiplication vector by itself
 */
func vectorItselfMultiplication(vector Vector) (Matrix) {
	lenVector := len(vector);
	out := Matrix{};

	for i := 0; i < lenVector; i++ {
		row := Vector{};
		for j := 0; j < lenVector; j++ {
			row = append(row, vector[j] * vector[i]);
		}

		out = append(out, row);
	}

	return out;
}

/**
Sum of matrices
 */
func sumMatrices(matrices []Matrix) (Matrix) {
	lenMatrix := len(matrices[0]);
	countMatrices := len(matrices);
	out := Matrix{};

	for i := 0; i< lenMatrix; i++ {
		row := Vector{};
		for j := 0; j< lenMatrix; j++ {
			var sum int64;
			for k := 0; k < countMatrices; k++ {
				sum += matrices[k][i][j];
			}

			row = append(row, sum);
		}

		out = append(out, row);
	}

	return out;
}

/**
Network convergence
 */
func convergence(matrix Matrix, vector Vector) (Vector) {
	lenMatrix := len(matrix[0]);
	out := Vector{};

	for i := 0; i< lenMatrix ; i++ {
		var sum int64;
		for j := 0; j< lenMatrix ; j++ {
			sum += matrix[i][j] * vector[j];
		}

		// Activation function
		if (sum >= 0) {
			sum = 1;
		} else {
			sum = -1;
		}

		out = append(out, sum);
	}

	return out;
}

/**
Get index sample
 */
func getIndexSample(vectors []Vector, vector Vector) (int) {
	for index, v := range vectors {
		if (reflect.DeepEqual(v, vector)) {
			return index;
		}
	};

	return -1;
}

/**
Write False Attractor
 */
func writeFalseAttractor(vector Vector) {
	f, _ := os.Create(DIR_SAMPLES + "/false_attractor");
	for index, char := range vector {
		if char >= 0 {
			f.WriteString("x");
		} else {
			f.WriteString(".");
		}

		if index % 6 == 5 {
			f.WriteString("\n");
		}
	};
}

func main() {
	if len(os.Args) <= 1 {
		println("No input vector");
		return;
	}

	samples := []string{};
	files, _ := ioutil.ReadDir(DIR_SAMPLES_ORIGINAL)
	for _, file := range files {
		samples = append(samples, DIR_SAMPLES_ORIGINAL + "/" + file.Name());
	}

	vectors := []Vector{};
	matrices := []Matrix{};

	for _, sample := range samples {
		vector := makeVector(sample);
		vectors = append(vectors, vector);
		matrices = append(matrices, vectorItselfMultiplication(vector));
	};

	WeightMatrix := sumMatrices(matrices);

	// Clear diagonal of matrix
	lenMatrix := len(WeightMatrix);
	for index := 0; index < lenMatrix; index++ {
		WeightMatrix[index][index] = 0;
	}

	lastVector := makeVector(os.Args[1]);
	for step := 1; step <= MAX_CONVERGENCE_STEPS; step++ {
		vector := convergence(WeightMatrix, lastVector);
		if (reflect.DeepEqual(vector, lastVector)) {
			index := getIndexSample(vectors, vector);
			if (index >= 0) {
				println("Sample found: " + samples[index] + ", Steps:" + strconv.Itoa(step));
				return;
			}

			writeFalseAttractor(vector);

			println("False attractor detected");

			return;
		}

		lastVector = vector;
	}

	println("Sample not found");
}