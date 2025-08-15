package pos

import (
	"fmt"

	"github.com/t0l1k/eui/utils"
)

const (
	TurnX     = 'x'
	TurnO     = 'o'
	TurnEmpty = ' '
)

type Posititon struct {
	dim, size int
	turn      rune
	board     []rune
	cache     map[int]int
}

func NewPosititon(dim int) *Posititon {
	p := &Posititon{dim: dim, size: dim * dim}
	p.Reset()
	return p
}

func (p *Posititon) Reset() {
	p.turn = TurnX
	p.cache = make(map[int]int)
	p.board = make([]rune, 0)
	for i := 0; i < p.size; i++ {
		p.board = append(p.board, TurnEmpty)
	}
}

func (p *Posititon) Setup(turn rune, board []rune) *Posititon {
	p.turn = turn
	p.board = nil
	p.board = board
	return p
}

func (p *Posititon) GetBoard() []rune { return p.board }
func (p *Posititon) GetTurn() string  { return string(p.turn) }
func (p *Posititon) GetNextTurn() string {
	turn := TurnEmpty
	if p.turn == TurnX {
		turn = TurnO
	} else {
		turn = TurnX
	}
	return string(turn)
}

func (p *Posititon) Move(idx int) *Posititon {
	p.board[idx] = p.turn
	if p.turn == TurnX {
		p.turn = TurnO
	} else {
		p.turn = TurnX
	}
	return p
}

func (p *Posititon) UnMove(idx int) *Posititon {
	p.board[idx] = TurnEmpty
	if p.turn == TurnX {
		p.turn = TurnO
	} else {
		p.turn = TurnX
	}
	return p
}

func (p *Posititon) IsWinFor(turn rune) (isWin bool) {
	for i := 0; i < p.size; i += p.dim {
		isWin = isWin || p.lineMatch(turn, i, i+p.dim, 1)
	}
	for i := 0; i < p.dim; i++ {
		isWin = isWin || p.lineMatch(turn, i, p.size, p.dim)
	}
	isWin = isWin || p.lineMatch(turn, 0, p.size, p.dim+1)
	isWin = isWin || p.lineMatch(turn, p.dim-1, p.size-1, p.dim-1)
	return isWin
}

func (p *Posititon) lineMatch(turn rune, start, end, step int) (isWin bool) {
	for i := start; i < end; i += step {
		if p.board[i] != turn {
			return false
		}
	}
	return true
}

func (p *Posititon) PossibleMoves() (result utils.IntList) {
	result = utils.NewIntList()
	for i, v := range p.board {
		if v == TurnEmpty {
			result = result.Add(i)
		}
	}
	return result
}

func (p *Posititon) Blanks() (total int) {
	for _, v := range p.board {
		if v == TurnEmpty {
			total++
		}
	}
	return total
}

func (p *Posititon) Code() (value int) {
	for i := 0; i < p.size; i++ {
		value = value * 3
		switch p.board[i] {
		case TurnX:
			value += 1
		case TurnO:
			value += 2
		}
	}
	return value
}

func (p *Posititon) Minimax() int {
	key := p.Code()
	value := p.cache[key]
	if value != 0 {
		return value
	}
	if p.IsWinFor(TurnX) {
		return p.Blanks()
	}
	if p.IsWinFor(TurnO) {
		return -p.Blanks()
	}
	if p.Blanks() == 0 {
		return 0
	}
	list := utils.NewIntList()
	for _, idx := range p.PossibleMoves() {
		list = list.Add(p.Move(idx).Minimax())
		p.UnMove(idx)
	}
	if p.turn == TurnX {
		value = list.Max()
	} else {
		value = list.Min()
	}
	p.cache[key] = value
	return value
}

func (p *Posititon) BestMove() int {
	init := false
	cur := 0
	cur_i := -1
	for i, v := range p.board {
		if v != TurnEmpty {
			continue
		}
		value := p.Move(i).Minimax()
		p.UnMove(i)
		if !init {
			init = true
			cur = value
			cur_i = i
		} else {
			if p.turn == TurnX && cur < value {
				cur = value
				cur_i = i
			} else if p.turn == TurnO && value < cur {
				cur = value
				cur_i = i
			}
		}
	}
	return cur_i
}

func (p *Posititon) IsGameEnd() bool {
	return p.IsWinFor(TurnX) || p.IsWinFor(TurnO) || p.Blanks() == 0
}

func (p Posititon) String() (result string) {
	result = fmt.Sprintf("TTT dim(%v) turn(%v) ", p.dim, string(p.turn))
	if p.IsGameEnd() {
		if p.IsWinFor(TurnX) {
			result += "Won X!"
		} else if p.IsWinFor(TurnO) {
			result += "Won O!"
		} else {
			result += "Draw!"
		}
	}
	result += "\n"
	for y := 0; y < p.dim; y++ {
		for x := 0; x < p.dim; x++ {
			result += fmt.Sprintf("[%v]", string(p.board[y*p.dim+x]))
		}
		result += "\n"
	}
	return result
}

func (p Posititon) StringShort() (result string) {
	for y := 0; y < p.dim; y++ {
		for x := 0; x < p.dim; x++ {
			result += fmt.Sprintf("%v", string(p.board[y*p.dim+x]))
		}
	}
	return result
}
