package state

type GameLogic interface {
	GameOver()
}

type LocalGameLogic struct {
}

func (l *LocalGameLogic) GameOver() {

}

type NetGameLogic struct {
}

func (l *NetGameLogic) GameOver() {
	
}
