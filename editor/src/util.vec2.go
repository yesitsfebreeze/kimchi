package main

type Vec2 struct {
	X int
	Y int
}

func (a Vec2) Add(b Vec2) Vec2 {
	return Vec2{a.X + b.X, a.Y + b.Y}
}

func (a Vec2) Sub(b Vec2) Vec2 {
	return Vec2{a.X - b.X, a.Y - b.Y}
}

func (a Vec2) Mul(b Vec2) Vec2 {
	return Vec2{a.X * b.X, a.Y * b.Y}
}

func (a Vec2) Div(b Vec2) Vec2 {
	return Vec2{
		X: func() int {
			if b.X != 0 {
				return a.X / b.X
			} else {
				return 0
			}
		}(),
		Y: func() int {
			if b.Y != 0 {
				return a.Y / b.Y
			} else {
				return 0
			}
		}(),
	}
}

func (a Vec2) Scale(s int) Vec2 {
	return Vec2{a.X * s, a.Y * s}
}
