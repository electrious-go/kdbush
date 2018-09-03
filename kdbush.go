package kdbush

// Interface, that should be implemented by indexing structure
// It's just simply returns points coordinates
// Called once, only when index created, so you could calc values on the fly for this interface
type Point interface {
	Coordinates() (X, Y float64)
}

// SimplePoint minimal struct, that implements Point interface
type SimplePoint struct {
	X, Y float64
}

// Coordinates to make SimplePoint's  implementation of Point interface satisfied
func (sp *SimplePoint) Coordinates() (float64, float64) {
	return sp.X, sp.Y
}

// KDBush a very fast static spatial index for 2D points based on a flat KD-tree.
// Points only, no rectangles
// static (no add, remove items)
// 2 dimensional
// indexing 16-40 times faster then  rtreego(https://github.com/dhconnelly/rtreego) (TODO: benchmark)
type KDBush struct {
	NodeSize int
	Points   []Point

	idxs   []int     //array of indexes
	coords []float64 //array of coordinates
}

// NewBush create new index from points
// Structure don't copy points itself, copy only coordinates
// Returns pointer to new KDBush index object, all data in it already indexed
// Input:
// points - slice of objects, that implements Point interface
// nodeSize  - size of the KD-tree node, 64 by default. Higher means faster indexing but slower search, and vise versa.
func NewBush(points []Point, nodeSize int) *KDBush {
	b := KDBush{}
	b.buildIndex(points, nodeSize)
	return &b
}

// Range finds all items within the given bounding box and returns an array of indices that refer to the items in the original points input slice.
func (bush *KDBush) Range(minX, minY, maxX, maxY float64) []int {
	stack := []int{0, len(bush.idxs) - 1, 0}
	result := []int{}
	var x, y float64
	for len(stack) > 0 {
		axis := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		right := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		left := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if right-left <= bush.NodeSize {
			for i := left; i <= right; i++ {
				x = bush.coords[2*i]
				y = bush.coords[2*i+1]
				if x >= minX && x <= maxX && y >= minY && y <= maxY {
					result = append(result, bush.idxs[i])
				}
			}
			continue
		}
		m := floor(float64(left+right) / 2.0)
		x = bush.coords[2*m]
		y = bush.coords[2*m+1]
		if x >= minX && x <= maxX && y >= minY && y <= maxY {
			result = append(result, bush.idxs[m])
		}
		nextAxis := (axis + 1) % 2
		if (axis == 0 && minX <= x) || (axis != 0 && minY <= y) {
			stack = append(stack, left)
			stack = append(stack, m-1)
			stack = append(stack, nextAxis)
		}
		if (axis == 0 && maxX >= x) || (axis != 0 && maxY >= y) {
			stack = append(stack, m+1)
			stack = append(stack, right)
			stack = append(stack, nextAxis)
		}
	}
	return result
}

// Within finds all items within a given radius from the query point and returns an array of indices.
func (bush *KDBush) Within(point Point, radius float64) []int {
	stack := []int{0, len(bush.idxs) - 1, 0}
	result := []int{}
	r2 := radius * radius
	qx, qy := point.Coordinates()
	for len(stack) > 0 {
		axis := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		right := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		left := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if right-left <= bush.NodeSize {
			for i := left; i <= right; i++ {
				dst := sqrtDist(bush.coords[2*i], bush.coords[2*i+1], qx, qy)
				if dst <= r2 {
					result = append(result, bush.idxs[i])
				}
			}
			continue
		}
		m := floor(float64(left+right) / 2.0)
		x := bush.coords[2*m]
		y := bush.coords[2*m+1]
		if sqrtDist(x, y, qx, qy) <= r2 {
			result = append(result, bush.idxs[m])
		}
		nextAxis := (axis + 1) % 2
		if (axis == 0 && (qx-radius <= x)) || (axis != 0 && (qy-radius <= y)) {
			stack = append(stack, left)
			stack = append(stack, m-1)
			stack = append(stack, nextAxis)
		}
		if (axis == 0 && (qx+radius >= x)) || (axis != 0 && (qy+radius >= y)) {
			stack = append(stack, m+1)
			stack = append(stack, right)
			stack = append(stack, nextAxis)
		}
	}
	return result
}

func (bush *KDBush) buildIndex(points []Point, nodeSize int) {
	bush.NodeSize = nodeSize
	bush.Points = points
	bush.idxs = make([]int, len(points))
	bush.coords = make([]float64, 2*len(points))
	for i, v := range points {
		bush.idxs[i] = i
		x, y := v.Coordinates()
		bush.coords[i*2] = x
		bush.coords[i*2+1] = y
	}
	sort(bush.idxs, bush.coords, bush.NodeSize, 0, len(bush.idxs)-1, 0)
}
