package geom

func NewRay(origin Tuple4, direction Tuple4) Ray {
	return Ray{Origin: origin, Direction: direction}
}
func NewEmptyRay() Ray {
	return Ray{Origin: NewTuple(), Direction: NewTuple()}
}

type Ray struct {
	Origin    Tuple4
	Direction Tuple4
}
