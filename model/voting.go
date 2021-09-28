package model

import (
	"math"

	"github.com/gitran-com/gitran-server/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Voting struct {
	UserID int64 `gorm:"primaryKey"`
	TranID int64 `gorm:"primaryKey"`
	Vote   int   `gorm:"vote"`
}

const (
	VoteLike   = 1
	VoteUnlike = -1
	VoteNone   = 0
)

//TableName returns table name
func (*Voting) TableName() string {
	return config.DB.TablePrefix + "votings"
}

func GetVote(user_id int64, tran_id int64) *Voting {
	var vote Voting
	res := db.Where("user_id=? AND tran_id=?", user_id, tran_id).First(&vote)
	if res.Error != nil {
		return nil
	}
	return &vote
}

func NewVote(user_id int64, tran_id int64, vote int) {
	db.Create(&Voting{
		UserID: user_id,
		TranID: tran_id,
		Vote:   vote,
	})
}

func UpdateVote(user_id int64, tran_id int64, vote int) {
	db.Raw("UPDATE votings SET vote=? WHERE user_id=? AND tran_id=?", vote, user_id, tran_id)
}

func calcLikesDelta(x, y int) int {
	if y == VoteLike && x != VoteLike {
		return 1
	}
	if x == VoteLike && y != VoteLike {
		return -1
	}
	return 0
}

func calcUnlikesDelta(x, y int) int {
	if y == VoteUnlike && x != VoteUnlike {
		return 1
	}
	if x == VoteUnlike && y != VoteUnlike {
		return -1
	}
	return 0
}

func wilsonScore(pos int, neg int) (score float64) {
	return func(pos int, total int, pz float64) float64 {
		if total == 0 {
			return 0.0
		}
		//https://en.wikipedia.org/wiki/Binomial_proportion_confidence_interval
		pos_rat := float64(pos) / float64(total)
		pz2 := pz * pz
		score := (pos_rat + (pz2 / (2.0 * float64(total))) - ((pz / (2.0 * float64(total))) * math.Sqrt(4.0*float64(total)*(1.0-pos_rat)*pos_rat+pz2))) / (1.0 + pz2/float64(total))
		return score
	}(pos, pos+neg, 2.0)
}

func TxnVote(user_id int64, tran_id int64, vote int) (likes int) {
	db.Transaction(func(tx *gorm.DB) error {
		var (
			oldVote = Voting{
				UserID: user_id,
				TranID: tran_id,
				Vote:   vote,
			}
			tran                       Translation
			sent                       Sentence
			likes_delta, unlikes_delta int
		)
		if res := tx.First(&tran, tran_id); res.Error != nil {
			return res.Error
		}
		likes = int(tran.Likes)
		if res := tx.First(&sent, tran.SentID); res.Error != nil {
			return res.Error
		}
		if res := tx.First(&oldVote); res.Error != nil {
			if res := tx.Create(&oldVote); res.Error != nil {
				return res.Error
			}
			if vote == VoteLike {
				likes_delta = 1
			} else if vote == VoteUnlike {
				unlikes_delta = 1
			}
		} else {
			if oldVote.Vote != vote {
				likes_delta = calcLikesDelta(oldVote.Vote, vote)
				unlikes_delta = calcUnlikesDelta(oldVote.Vote, vote)
				if res := tx.Model(&oldVote).Update("vote", vote); res.Error != nil {
					return res.Error
				}
			}
		}
		if likes_delta != 0 || unlikes_delta != 0 {
			score := wilsonScore(likes+likes_delta, int(tran.Unlikes)+unlikes_delta)
			// fmt.Printf("score(%v,%v)=%v\n", likes+likes_delta, int(tran.Unlikes)+unlikes_delta, score)
			res := tx.Model(&tran).Updates(map[string]interface{}{"likes": gorm.Expr("likes+(?)", likes_delta), "unlikes": gorm.Expr("unlikes+(?)", unlikes_delta), "score": score})
			if res.Error != nil {
				log.Warnf("TxnVote UPDATE ERROR: %v", res.Error.Error())
				return res.Error
			}
		}
		likes = int(tran.Likes) + likes_delta
		return nil
	})
	return
}
