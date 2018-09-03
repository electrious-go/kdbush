package kdbush

import "math"

func sort(idxs []int, coords []float64, nodeSize int, left, right, depth int) {
	if (right - left) <= nodeSize {
		return
	}
	m := floor(float64(left+right) / 2.0)
	sselect(idxs, coords, m, left, right, depth%2)
	sort(idxs, coords, nodeSize, left, m-1, depth+1)
	sort(idxs, coords, nodeSize, m+1, right, depth+1)

}

func sselect(idxs []int, coords []float64, k, left, right, inc int) {
	//whatever you want
	for right > left {
		if (right - left) > 600 {
			n := right - left + 1
			m := k - left + 1
			z := math.Log(float64(n))
			s := 0.5 * math.Exp(2.0*z/3.0)
			sds := 1.0
			if float64(m)-float64(n)/2.0 < 0 {
				sds = -1.0
			}
			ns := float64(n) - s
			sd := 0.5 * math.Sqrt(z*s*ns/float64(n)) * sds
			newLeft := iMax(left, floor(float64(k)-float64(m)*s/float64(n)+sd))
			newRight := iMin(right, floor(float64(k)+float64(n-m)*s/float64(n)+sd))
			sselect(idxs, coords, k, newLeft, newRight, inc)
		}

		t := coords[2*k+inc]
		i := left
		j := right

		swapItem(idxs, coords, left, k)
		if coords[2*right+inc] > t {
			swapItem(idxs, coords, left, right)
		}

		for i < j {
			swapItem(idxs, coords, i, j)
			i++
			j--
			for coords[2*i+inc] < t {
				i++
			}
			for coords[2*j+inc] > t {
				j--
			}
		}

		if coords[2*left+inc] == t {
			swapItem(idxs, coords, left, j)
		} else {
			j++
			swapItem(idxs, coords, j, right)
		}

		if j <= k {
			left = j + 1
		}
		if k <= j {
			right = j - 1
		}
	}
}

func swapItem(idxs []int, coords []float64, i, j int) {
	swapi(idxs, i, j)
	swapf(coords, 2*i, 2*j)
	swapf(coords, 2*i+1, 2*j+1)
}

func swapf(a []float64, i, j int) {
	t := a[i]
	a[i] = a[j]
	a[j] = t
}

func swapi(a []int, i, j int) {
	t := a[i]
	a[i] = a[j]
	a[j] = t
}

func iMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func iMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func floor(in float64) int {
	out := math.Floor(in)
	return int(out)
}

func sqrtDist(ax, ay, bx, by float64) float64 {
	dx := ax - bx
	dy := ay - by
	return dx*dx + dy*dy
}
