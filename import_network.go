package cnns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/LdDl/cnns/tensor"
	"gonum.org/v1/gonum/mat"
)

// NetJSON JSON representation of network structure (for import and export)
type NetJSON struct {
	Network    *NetworkJSON    `json:"network"`
	Parameters *LearningParams `json:"parameters"`
}

// NetworkJSON JSON representation of networks' layers
type NetworkJSON struct {
	Layers []*NetLayerJSON `json:"layers"`
}

// NetLayerJSON JSON representation of layer
type NetLayerJSON struct {
	LayerType  string           `json:"layer_type"`
	InputSize  *tensor.TDsize   `json:"input_size"`
	Parameters *LayerParamsJSON `json:"parameters"`
	Weights    []*NestedData    `json:"weights"`
	// Actually "OutputSize" parameter is useful for fully-connected layer only
	// There are automatic calculation of output size for other layers' types
	OutputSize *tensor.TDsize `json:"output_size"`
}

// LayerParamsJSON JSON representation of layers attributes
type LayerParamsJSON struct {
	Stride          int    `json:"stride"`
	KernelSize      int    `json:"kernel_size"`
	PoolingType     string `json:"pooling_type"`
	ZeroPaddingType string `json:"zero_padding_type"`
}

// NestedData JSON representation of stored data
type NestedData struct {
	Data []float64 `json:"data"`
}

// ImportFromFile Load network to file
/*
	fname - filename,
	randomWeights:
		true: random weights for new network
		false: weights from files for using network (or continue training))
*/
func (wh *WholeNet) ImportFromFile(fname string, randomWeights bool) error {
	var err error
	fileBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	var data NetJSON
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return err
	}
	// @todo Need to handle case when wh.Layers is not empty slice
	wh.Layers = []Layer{}
	for i := range data.Network.Layers {
		switch data.Network.Layers[i].LayerType {
		case "conv":
			stride := data.Network.Layers[i].Parameters.Stride
			kernelSize := data.Network.Layers[i].Parameters.KernelSize
			numOfFilters := len(data.Network.Layers[i].Weights)
			x := data.Network.Layers[i].InputSize.X
			y := data.Network.Layers[i].InputSize.Y
			z := data.Network.Layers[i].InputSize.Z
			conv := NewConvLayer(&tensor.TDsize{X: x, Y: y, Z: z}, stride, kernelSize, numOfFilters)
			if randomWeights == false {
				weights := make([]*mat.Dense, numOfFilters)
				for w := 0; w < numOfFilters; w++ {
					weights[w] = mat.NewDense(kernelSize*z, kernelSize, data.Network.Layers[i].Weights[w].Data)
				}
				conv.SetCustomWeights(weights)
			}
			wh.Layers = append(wh.Layers, conv)
			break
		case "relu":
			x := data.Network.Layers[i].InputSize.X
			y := data.Network.Layers[i].InputSize.Y
			z := data.Network.Layers[i].InputSize.Z
			relu := NewReLULayer(&tensor.TDsize{X: x, Y: y, Z: z})
			wh.Layers = append(wh.Layers, relu)
			break
		case "pool":
			stride := data.Network.Layers[i].Parameters.Stride
			kernelSize := data.Network.Layers[i].Parameters.KernelSize
			x := data.Network.Layers[i].InputSize.X
			y := data.Network.Layers[i].InputSize.Y
			z := data.Network.Layers[i].InputSize.Z
			pool := NewPoolingLayer(&tensor.TDsize{X: x, Y: y, Z: z}, stride, kernelSize, data.Network.Layers[i].Parameters.PoolingType, data.Network.Layers[i].Parameters.ZeroPaddingType)
			wh.Layers = append(wh.Layers, pool)
			break
		case "fc":
			x := data.Network.Layers[i].InputSize.X
			y := data.Network.Layers[i].InputSize.Y
			z := data.Network.Layers[i].InputSize.Z
			outSize := data.Network.Layers[i].OutputSize.X
			fullyconnected := NewFullyConnectedLayer(&tensor.TDsize{X: x, Y: y, Z: z}, outSize)
			if randomWeights == false {
				weights := mat.NewDense(outSize, x*y*z, data.Network.Layers[i].Weights[0].Data)
				fullyconnected.SetCustomWeights([]*mat.Dense{weights})
			}
			wh.Layers = append(wh.Layers, fullyconnected)
			break
		default:
			err = fmt.Errorf("Unrecognized layer type: %s", data.Network.Layers[i].LayerType)
			return err
		}
	}

	wh.LP.LearningRate = data.Parameters.LearningRate
	wh.LP.Momentum = data.Parameters.Momentum

	return err
}
