package game

type Player struct {
	Paddle *Paddle
}

func NewPlayer(paddle *Paddle) *Player {
	return &Player{
		Paddle: paddle,
	}
}
