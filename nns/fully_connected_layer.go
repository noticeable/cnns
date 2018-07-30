package nns

import (
	"fmt"
	"math"
	"math/rand"
)

// FullConnectedLayer is simple layer structure (so this layer can be used for simple neural networks like XOR problem), where
type FullConnectedLayer struct {
	In                    Tensor
	Out                   Tensor
	InputGradientsWeights Tensor
	Weights               Tensor
	LocalGradients        []Gradient
	Input                 []float64
	ActivationFunc        func(v float64) float64
	ActivationDerivative  func(v float64) float64
}

// ActivationSygmoid is default ActivationFunc
func ActivationSygmoid(v float64) float64 {
	return 1.0 / (1.0 + math.Exp(-1.0*v))
}

// ActivationSygmoidDerivative is default derivative of ActivationFunc (which is ActivationSygmoid)
func ActivationSygmoidDerivative(v float64) float64 {
	return ActivationSygmoid(v) * (1 - ActivationSygmoid(v))
}

// SetActivationFunc sets activation function for fully connected layer. You need to specify function: func(v float64) float64
func (fc *FullConnectedLayer) SetActivationFunc(f func(v float64) float64) {
	(*fc).ActivationFunc = f
}

// SetActivationDerivativeFunc sets derivative of activation function for fully connected layer. You need to specify function: func(v float64) float64
func (fc *FullConnectedLayer) SetActivationDerivativeFunc(f func(v float64) float64) {
	(*fc).ActivationDerivative = f
}

// NewFullConnectedLayer - constructor for new fully connected layer. You need to specify input size and output size
func NewFullConnectedLayer(inSize TDsize, outSize int) *LayerStruct {
	newLayer := &FullConnectedLayer{
		In:  NewTensor(inSize.X, inSize.Y, inSize.Z),
		Out: NewTensor(outSize, 1, 1),
		InputGradientsWeights: NewTensor(inSize.X, inSize.Y, inSize.Z),
		Weights:               NewTensor(inSize.X*inSize.Y*inSize.Z, outSize, 1),
		Input:                 make([]float64, outSize),
		LocalGradients:        make([]Gradient, outSize),
		ActivationFunc:        ActivationSygmoid,           // Default Activation function is Sygmoid
		ActivationDerivative:  ActivationSygmoidDerivative, // Default derivative of activation function is Sygmoid*(1-Sygmoid)
	}
	for i := 0; i < outSize; i++ {
		for h := 0; h < inSize.X*inSize.Y*inSize.Z; h++ {
			newLayer.Weights.Set(h, i, 0, rand.Float64()-0.5)
		}
	}
	// hardcoded weights for testing purposes
	// hcweights := [][][]float64{
	// 	[][]float64{
	// 		[]float64{-0.19908814, 0.01521263, 0.31363996, -0.28573613, -0.11934281, -0.18194183, -0.03111016, -0.21696585, -0.20689814},
	// 		[]float64{0.17908468, -0.28144695, -0.29681312, -0.13912858, 0.07067328, 0.36249144, -0.20688576, -0.20291744, 0.25257304},
	// 		[]float64{-0.29341734, 0.36533501, 0.19671917, 0.02382031, -0.47169692, -0.34167172, 0.10725344, 0.47524162, -0.42054638},
	// 	},
	// }
	// newLayer.Weights.CopyFrom(hcweights)
	return &LayerStruct{
		Layer: newLayer,
	}
}

func (fc *FullConnectedLayer) OutSize() Point {
	return (*fc).Out.Size
}

func (fc *FullConnectedLayer) mapToInput(d Point) int {
	return d.Z*(*fc).In.Size.X*(fc).In.Size.Y + d.Y*(*fc).In.Size.X + d.X
}

// PrintWeights - print fully connected layer's weights
func (fc *FullConnectedLayer) PrintWeights() {
	fmt.Println("Printing Fully Connected Layer weights...")
	(*fc).Weights.Print()
}

// PrintOutput - print fully connected layer's output
func (fc *FullConnectedLayer) PrintOutput() {
	fmt.Println("Printing Fully Connected Layer output...")
	(*fc).Out.Print()
}

// PrintGradients - print fully connected layer's gradients
func (fc *FullConnectedLayer) PrintGradients() {
	fmt.Println("Printing Fully Connected Layer gradients-weights...")
	(*fc).InputGradientsWeights.Print()
}

// PrintSumGradWeights - print fully connected layer's summ of grad*weight
func (fc *FullConnectedLayer) PrintSumGradWeights() {
	fmt.Println("Printing Fully Connected Layer summ of grad*weight...")
	//(*fc).SumDeltaWeights.Print()
}

// GetOutput - get fully connected layer's output
func (fc *FullConnectedLayer) GetOutput() Tensor {
	return (*fc).Out // Here we outputing ACTIVATED values
}

// FeedForward - feed data to fully connected layer
func (fc *FullConnectedLayer) FeedForward(t *Tensor) {
	(*fc).In = (*t)
	(*fc).DoActivation()
}

// GetGradients - get sum (next layer grad * weights)
func (fc *FullConnectedLayer) GetGradients() Tensor {
	return (*fc).InputGradientsWeights
}

// CalculateGradients - calculate fully connected layer's gradients
func (fc *FullConnectedLayer) CalculateGradients(nextLayerGradients *Tensor) {
	for i := 0; i < (*fc).InputGradientsWeights.Size.X*(*fc).InputGradientsWeights.Size.Y*(*fc).InputGradientsWeights.Size.Z; i++ {
		(*fc).InputGradientsWeights.Data[i] = 0.0
	}
	for n := 0; n < (*fc).Out.Size.X; n++ {
		(*fc).LocalGradients[n].Grad = (*nextLayerGradients).Get(n, 0, 0) * (*fc).ActivationDerivative((*fc).Input[n])
		for i := 0; i < (*fc).In.Size.X; i++ {
			for j := 0; j < (*fc).In.Size.Y; j++ {
				for z := 0; z < (*fc).In.Size.Z; z++ {
					m := fc.mapToInput(Point{X: i, Y: j, Z: z})
					// fmt.Printf("%v * %v\n", (*fc).LocalGradients[n].Grad, (*fc).Weights.Get(m, n, 0))
					(*fc).InputGradientsWeights.SetAdd(i, j, z, (*fc).LocalGradients[n].Grad*(*fc).Weights.Get(m, n, 0))
				}
			}
		}
	}
}

// UpdateWeights - update fully connected layer's weights
func (fc *FullConnectedLayer) UpdateWeights() {
	for n := 0; n < (*fc).Out.Size.X; n++ {
		grad := (*fc).LocalGradients[n]
		for i := 0; i < (*fc).In.Size.X; i++ {
			for j := 0; j < (*fc).In.Size.Y; j++ {
				for z := 0; z < (*fc).In.Size.Z; z++ {
					m := fc.mapToInput(Point{X: i, Y: j, Z: z})
					w := (*fc).Weights.Get(m, n, 0)
					w = UpdateWeight(w, &grad, (*fc).In.Get(i, j, z))
					(*fc).Weights.Set(m, n, 0, w)
				}
			}
		}
		UpdateGradient(&fc.LocalGradients[n])
	}
}

// DoActivation - fully connected layer's output activation
func (fc *FullConnectedLayer) DoActivation() {
	for n := 0; n < (*fc).Out.Size.X; n++ {
		inputv := 0.0
		for i := 0; i < (*fc).In.Size.X; i++ {
			for j := 0; j < (*fc).In.Size.Y; j++ {
				for z := 0; z < (*fc).In.Size.Z; z++ {
					m := fc.mapToInput(Point{X: i, Y: j, Z: z})
					inputv += (*fc).In.Get(i, j, z) * (*fc).Weights.Get(m, n, 0)
				}
			}
		}
		(*fc).Input[n] = inputv
		(*fc).Out.Set(n, 0, 0, (*fc).ActivationFunc(inputv))
	}
}
