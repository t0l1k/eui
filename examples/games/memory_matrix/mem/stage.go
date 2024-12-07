package mem

import (
	"log"
)

type GameStage int

const (
	Preparation GameStage = iota
	Show
	Recollection
	Conclusion
	Restart
)

func NewGameStage() GameStage            { return Preparation }
func (g GameStage) IsPreparation() bool  { return g == Preparation }
func (g GameStage) IsShow() bool         { return g == Show }
func (g GameStage) IsRecollection() bool { return g == Recollection }
func (g GameStage) IsConclusion() bool   { return g == Conclusion }
func (g GameStage) IsRestart() bool      { return g == Restart }
func (g GameStage) Next() GameStage {
	switch g {
	case Preparation:
		g = Show
	case Show:
		g = Recollection
	case Recollection:
		g = Conclusion
	case Conclusion:
		g = Restart
	case Restart:
		g = Preparation
	}
	log.Println("set next srate:", g)
	return g
}

func (g GameStage) String() (result string) {
	switch g {
	case Preparation:
		result = "Preparation"
	case Show:
		result = "Show"
	case Recollection:
		result = "Recollection"
	case Conclusion:
		result = "Conclusion"
	case Restart:
		result = "Restart"
	}
	return result
}
