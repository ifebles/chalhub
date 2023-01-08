package checkerbot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ifebles/chalhub/pkg/util"
)

type playerType string

const (
	Ai    playerType = "ai"
	Human playerType = "human"
)

type PlayMode uint8

const (
	PlayerVsAI PlayMode = iota
	PlayerVsPlayer
	AIvsAI
)

func (mode PlayMode) String() string {
	switch mode {
	case PlayerVsAI:
		return "Player vs AI"

	case PlayerVsPlayer:
		return "Player vs Player"

	case AIvsAI:
		return "AI vs AI"
	}

	return "unknown"
}

var PlayModes = [3]PlayMode{
	PlayerVsAI, PlayerVsPlayer, AIvsAI,
}

type player struct {
	Char   rune
	Enemy  *player
	Type   playerType
	Pieces []Piece
}

func (pl player) containsPieceAt(p point) (bool, error) {
	if p.X < 0 || p.Y < 0 || p.X >= boardSize || p.Y >= boardSize {
		return false, fmt.Errorf("point out of bounds")
	}

	_, ok := util.Find(pl.Pieces, func(i Piece) bool {
		return i.Point.X == p.X && i.Point.Y == p.Y
	})

	return ok, nil
}

type movement struct {
	from, to point
	slay     *point
	isKing   bool
}

type moveParams struct {
	cond1, cond2 bool
	dir1, dir2   int
}

var currentTurn = 1
var players = [2]player{}

func StartGame(mode PlayMode) player {
	if !Board.initialized {
		go Board.initialize()
	}

	switch mode {
	case PlayerVsAI:
		rand.Seed(time.Now().UnixNano())

		if rand.Intn(2) == 1 {
			players[0] = player{blackChar, &players[1], Human, Board.black}
			players[1] = player{whiteChar, &players[0], Ai, Board.white}
		} else {
			players[0] = player{blackChar, &players[1], Ai, Board.black}
			players[1] = player{whiteChar, &players[0], Human, Board.white}
		}

	case PlayerVsPlayer:
		players[0] = player{blackChar, &players[1], Human, Board.black}
		players[1] = player{whiteChar, &players[0], Human, Board.white}

	case AIvsAI:
		players[0] = player{blackChar, &players[1], Ai, Board.black}
		players[1] = player{whiteChar, &players[0], Ai, Board.white}

	default:
		panic("unknown mode")
	}

	return players[0]
}

func GetTurnNumber() int {
	return currentTurn
}

func GetCurrentPlayer() player {
	return players[(currentTurn-1)%len(players)]
}

func EndTurn() player {
	currentTurn++

	return GetCurrentPlayer()
}

func PointToCoord(p point) string {
	if p.X < 0 || p.X >= boardSize || p.Y < 0 || p.Y >= boardSize {
		panic(fmt.Sprintf("invalid point: %v", p))
	}

	literal, numeral := rune('A'+p.Y), boardSize-p.X

	return fmt.Sprintf("%s%d", string(literal), numeral)
}

func CoordToPoint(c string) point {
	if len(c) != 2 {
		panic(fmt.Sprintf("invalid coordinate: %q", c))
	}

	literal := rune(strings.ToUpper(c)[0])

	if literal < 'A' || literal > 'A'+boardSize-1 {
		panic(fmt.Sprintf("invalid literal: %s", string(literal)))
	}

	numeral, err := strconv.Atoi(string(c[1]))

	if err != nil {
		panic(err)
	}

	if numeral <= 0 || numeral > boardSize {
		panic(fmt.Sprintf("invalid numeral: %d", numeral))
	}

	return point{boardSize - numeral, int(literal - 'A')}
}

func identifyMoves(pl player, po point, chkoverlap func(point) bool, dir vdirection) []movement {
	checkMv := func(xcond, ycond bool, xdir, ydir int) (movement, bool) {
		var mt movement

		if xcond && ycond {
			p := point{po.X + xdir, po.Y + ydir}
			isEnemy, err := pl.Enemy.containsPieceAt(p)

			if err != nil {
				return mt, false
			}

			if t, _ := pl.containsPieceAt(p); !isEnemy && !t {
				king := dir == vboth

				if !king {
					// Set king status for the movement
					king = (dir == up && p.X == 0) || (dir == down && p.X == boardSize-1)
				}

				return movement{po, p, nil, king}, true
			}

			if isEnemy {
				s := &point{p.X, p.Y}
				p = point{p.X + xdir, p.Y + ydir}
				king := dir == vboth

				if !king {
					// Set king status for the movement
					king = (dir == up && p.X == 0) || (dir == down && p.X == boardSize-1)
				}

				overlap := false

				if chkoverlap != nil {
					overlap = chkoverlap(p)
				}

				if _, ok, err := Board.GetPieceAt(p); err == nil && (overlap || !ok) {
					return movement{po, p, s, king}, true
				}
			}
		}

		return mt, false
	}

	mvs := []movement{}
	params := []moveParams{}

	switch dir {
	case up:
		params = append(
			params,
			// Top-left
			moveParams{po.X > 0, po.Y > 0, -1, -1},
			// Top-right
			moveParams{po.X > 0, po.Y <= boardSize, -1, 1},
		)

	case down:
		params = append(
			params,
			// Bottom-left
			moveParams{po.X <= boardSize, po.Y > 0, 1, -1},
			// Bottom-right
			moveParams{po.X <= boardSize, po.Y <= boardSize, 1, 1},
		)

	default: // both
		params = append(
			params,
			// Top-left
			moveParams{po.X > 0, po.Y > 0, -1, -1},
			// Top-right
			moveParams{po.X > 0, po.Y <= boardSize, -1, 1},
			// Bottom-left
			moveParams{po.X <= boardSize, po.Y > 0, 1, -1},
			// Bottom-right
			moveParams{po.X <= boardSize, po.Y <= boardSize, 1, 1},
		)
	}

	for _, a := range params {
		if mv, ok := checkMv(a.cond1, a.cond2, a.dir1, a.dir2); ok {
			mvs = append(mvs, mv)
		}
	}

	return mvs
}

func FilterSlayingOptions(pl player) []Piece {
	var dir vdirection

	if pl.Char == whiteChar {
		dir = up
	} else {
		dir = down
	}

	result := util.Filter(pl.Pieces, func(i Piece) bool {
		d := dir

		if i.IsKing {
			d = vboth
		}

		mvs := identifyMoves(pl, i.Point, nil, d)
		_, ok := util.Find(mvs, func(it movement) bool { return it.slay != nil })

		return ok
	})

	return result
}

func FilterSimpleOptions(pl player) []Piece {
	var dir vdirection

	if pl.Char == whiteChar {
		dir = up
	} else {
		dir = down
	}

	result := util.Filter(pl.Pieces, func(i Piece) bool {
		d := dir

		if i.IsKing {
			d = vboth
		}

		if mvs := identifyMoves(pl, i.Point, nil, d); len(mvs) > 0 {
			return true
		}

		return false
	})

	return result
}

func CreateTreeMaps(pl player, pi *Piece) []tree[movement] {
	result := []tree[movement]{}
	var dir vdirection

	if pi.IsKing {
		dir = vboth
	} else if pl.Char == whiteChar {
		dir = up
	} else {
		dir = down
	}

	mvs := identifyMoves(pl, pi.Point, nil, dir)

	if len(mvs) == 0 {
		panic("the piece has no valid move")
	}

	_, onlySlay := util.Find(mvs, func(i movement) bool { return i.slay != nil })

	for _, a := range mvs {
		if onlySlay && a.slay != nil {
			t := tree[movement]{dir, &xtreeNode[movement]{value: a}}
			populateTree(pl, t.start, a.from, dir, onlySlay)

			result = append(result, t)
		} else if !onlySlay {
			result = append(result, tree[movement]{dir, &xtreeNode[movement]{value: a}})
		}
	}

	return result
}

func populateTree(pl player, nd *xtreeNode[movement], st point, dir vdirection, slay bool) {
	var mvs []movement
	overlapchk := func(p point) bool {
		return p == st && p != nd.value.from
	}

	if m := identifyMoves(pl, nd.value.to, overlapchk, dir); slay {
		mvs = util.Filter(m, func(i movement) bool { return i.slay != nil })
	} else {
		mvs = m
	}

	if len(mvs) == 0 {
		return
	}

	for _, a := range mvs {
		n := &xtreeNode[movement]{value: a}
		var err error
		dirpair := struct {
			h hdirection
			v vdirection
		}{}

		if a.to.X < a.from.X {
			// Top-left
			if a.to.Y < a.from.Y {
				dirpair.v, dirpair.h = up, left
				err = nd.add(n, dirpair.v, dirpair.h)
			} else { // Top-right
				dirpair.v, dirpair.h = up, right
				err = nd.add(n, dirpair.v, dirpair.h)
			}
		} else {
			// Bottom-left
			if a.to.Y < a.from.Y {
				dirpair.v, dirpair.h = down, left
				err = nd.add(n, dirpair.v, dirpair.h)
			} else { // Bottom-right
				dirpair.v, dirpair.h = down, right
				err = nd.add(n, dirpair.v, dirpair.h)
			}
		}

		if err == nil {
			fnd := nd.findNode(func(v movement) bool { return a.to == v.from })

			if fnd != nil {
				var nuv vdirection

				if dirpair.v == up {
					nuv = down
				} else {
					nuv = up
				}

				// Connect both ends:
				n.set(fnd, nuv, dirpair.h)
			}

			////

			nudir := dir

			if nudir != vboth && a.isKing {
				// In case the move caused crowning:
				nudir = vboth
			}

			populateTree(pl, n, st, nudir, slay)
		}
	}
}
