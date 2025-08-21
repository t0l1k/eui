package mem

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
	return []GameStage{
		Show,
		Recollection,
		Conclusion,
		Restart,
		Preparation,
	}[g]
}

func (g GameStage) String() string {
	return []string{
		"Preparation",
		"Show",
		"Recollection",
		"Conclusion",
		"Restart",
	}[g]
}
