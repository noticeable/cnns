package tensor

type ConvLayer struct {
	In         *Tensor
	Out        *Tensor
	Kernels    []*Tensor
	KernelSize int
	Stride     int
}

type Tensor struct {
	Data []float64
	Size *Point
}

func NewTensor(x, y, z int) *Tensor {
	return &Tensor{
		Data: make([]float64, x*y*z),
		Size: &TDsize{
			X: x,
			Y: y,
			Z: z,
		},
	}
}
func (t *Tensor) Get(x, y, z int) float64 {
	return t.Data[z*t.Size.X*t.Size.Y+y*t.Size.X+x]
}
func (t *Tensor) Set(x, y, z int, val float64) {
	t.Data[z*t.Size.X*t.Size.Y+y*t.Size.X+x] = val
}

func NewConvolveLayer(insizeX, insizeY, insizeZ, numberFilters, kernelSize, stride int) *ConvLayer {
	tmp := &ConvLayer{
		In:         NewTensor(insizeX, insizeY, insizeZ),
		Out:        NewTensor((insizeX-kernelSize)/stride+1, (insizeY-kernelSize)/stride+1, numberFilters),
		KernelSize: kernelSize,
		Stride:     stride,
	}
	for i := 0; i < numberFilters; i++ {
		tmp.Kernels = append(tmp.Kernels, NewTensor(kernelSize, kernelSize, insizeZ))
	}
	return tmp
}

func (conv *ConvLayer) NaiveConv() {
	for filter := 0; filter < len(conv.Kernels); filter++ {
		filterData := conv.Kernels[filter]
		for x := 0; x < conv.Out.Size.X; x++ {
			for y := 0; y < conv.Out.Size.Y; y++ {
				mappedX, mappedY := x*conv.Stride, y*conv.Stride
				sum := 0.0
				for i := 0; i < conv.KernelSize; i++ {
					for j := 0; j < conv.KernelSize; j++ {
						for z := 0; z < conv.In.Size.Z; z++ {
							f := filterData.Get(i, j, z)
							v := conv.In.Get(mappedX+i, mappedY+j, z)
							sum += f * v
						}
					}
				}
				conv.Out.Set(x, y, filter, sum)
			}
		}
	}
}
