package game

type Player struct {
	Name   string
	Paddle *Paddle
}

func NewPlayer(name string, paddle *Paddle) *Player {
	return &Player{
		Name:   name,
		Paddle: paddle,
	}
}
