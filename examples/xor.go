package examples

import (
	"cnns_vika/nns"
	"cnns_vika/utils/u"
	"math/rand"
	"time"
)

// CheckXOR - проверка полносвязного слоя при решении проблемы XOR
func CheckXOR() {
	rand.Seed(time.Now().UnixNano())
	// Слой с тремя нейронами
	fullyconnected1 := nns.NewFullConnectedLayer(nns.TDsize{X: 2, Y: 1, Z: 1}, 2)
	fullyconnected1.SetActivationFunc(ActivationTanh)
	fullyconnected1.SetActivationDerivativeFunc(ActivationTanhDerivative)

	// Слой с одним выходным нейроном
	fullyconnected2 := nns.NewFullConnectedLayer(fullyconnected1.OutSize(), 1)
	fullyconnected2.SetActivationFunc(ActivationTanh)
	fullyconnected2.SetActivationDerivativeFunc(ActivationTanhDerivative)

	// Инициализация сети
	var net nns.WholeNet
	net.Layers = append(net.Layers, fullyconnected1)
	net.Layers = append(net.Layers, fullyconnected2)
	for i := 0; i < 100000; i++ {
		firstInt := u.RandomInt(0, 2)
		secondInt := u.RandomInt(0, 2)
		firstBool := false
		secondBool := false
		if firstInt == 1 {
			firstBool = true
		}
		if secondInt == 1 {
			secondBool = true
		}
		outputBool := (firstBool != secondBool)
		outputInt := 0
		if outputBool == true {
			outputInt = 1
		}
		desired := nns.NewTensor(1, 1, 1)
		desired.CopyFrom([][][]float64{[][]float64{[]float64{float64(outputInt)}}})
		input := nns.NewTensor(2, 1, 1)
		input.CopyFrom([][][]float64{[][]float64{[]float64{float64(firstInt), float64(secondInt)}}})
		// Forward
		net.Layers[0].FeedForward(&input)
		for l := 1; l < len(net.Layers); l++ {
			out := net.Layers[l-1].GetOutput()
			net.Layers[l].FeedForward(&out)
		}
		//Backward
		difference := net.Layers[len(net.Layers)-1].GetOutput()
		difference.Sub(&desired)

		net.Layers[len(net.Layers)-1].CalculateGradients(&difference)
		for i := len(net.Layers) - 2; i >= 0; i-- {
			grad := net.Layers[i+1].GetGradients()
			net.Layers[i].CalculateGradients(&grad)
		}
		for i := range net.Layers {
			net.Layers[i].UpdateWeights()
		}
	}

	// 0 * 0
	inputTest := nns.NewTensor(2, 1, 1)
	inputTest.CopyFrom([][][]float64{[][]float64{[]float64{0, 0}}})
	net.Layers[0].FeedForward(&inputTest)
	for l := 1; l < len(net.Layers); l++ {
		out := net.Layers[l-1].GetOutput()
		net.Layers[l].FeedForward(&out)
	}
	net.Layers[1].PrintOutput()

	// 1 * 0
	inputTest.CopyFrom([][][]float64{[][]float64{[]float64{1.0, 0}}})
	net.Layers[0].FeedForward(&inputTest)
	for l := 1; l < len(net.Layers); l++ {
		out := net.Layers[l-1].GetOutput()
		net.Layers[l].FeedForward(&out)
	}
	net.Layers[1].PrintOutput()

	// 0 * 1
	inputTest.CopyFrom([][][]float64{[][]float64{[]float64{0, 1.0}}})
	net.Layers[0].FeedForward(&inputTest)
	for l := 1; l < len(net.Layers); l++ {
		out := net.Layers[l-1].GetOutput()
		net.Layers[l].FeedForward(&out)
	}
	net.Layers[1].PrintOutput()

	// 1 * 1
	inputTest.CopyFrom([][][]float64{[][]float64{[]float64{1.0, 1.0}}})
	net.Layers[0].FeedForward(&inputTest)
	for l := 1; l < len(net.Layers); l++ {
		out := net.Layers[l-1].GetOutput()
		net.Layers[l].FeedForward(&out)
	}
	net.Layers[1].PrintOutput()
}
