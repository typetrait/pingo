package game

type Rules struct {
	WinningScore int64
}

func NewRules(winningScore int64) *Rules {
	return &Rules{
		WinningScore: winningScore,
	}
}
