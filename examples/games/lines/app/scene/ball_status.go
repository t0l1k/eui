package scene

import "strconv"

type BallStatusType int

const (
	BallHidden BallStatusType = iota
	BallSmall
	BallMedium
	BallNormal
	BallBig

	BallJumpDown
	BallJumpCenter
	BallJumpUp
)

func (j *BallStatusType) IsHidden() bool { return *j == BallHidden }
func (j *BallStatusType) IsSmall() bool  { return *j == BallSmall }
func (j *BallStatusType) IsMedium() bool { return *j == BallMedium }
func (j *BallStatusType) IsNormal() bool { return *j == BallNormal }
func (j *BallStatusType) IsBig() bool    { return *j == BallBig }

func (j *BallStatusType) FilledNext() {
	switch *j {
	case BallHidden:
		*j = BallSmall
	case BallSmall:
		*j = BallMedium
	}
}

func (j *BallStatusType) Jump() {
	switch *j {
	case BallJumpDown:
		*j = BallJumpCenter
	case BallJumpCenter:
		*j = BallJumpUp
	case BallJumpUp:
		*j = BallJumpDown
	}
}

func (j *BallStatusType) Delete() {
	switch *j {
	case BallBig:
		*j = BallNormal
	case BallNormal:
		*j = BallMedium
	case BallMedium:
		*j = BallSmall
	case BallSmall:
		*j = BallHidden
	}
}

func (j BallStatusType) String() (res string) {
	switch j {
	case BallHidden:
		res = "ball hidden"
	case BallSmall:
		res = "ball small"
	case BallMedium:
		res = "ball mediun"
	case BallNormal:
		res = "ball normal"
	case BallBig:
		res = "ball big"
	case BallJumpDown:
		res = "jump down"
	case BallJumpCenter:
		res = "jump center"
	case BallJumpUp:
		res = "jump up"
	default:
		res = strconv.Itoa(int(j)) + "!"
	}
	return res
}
