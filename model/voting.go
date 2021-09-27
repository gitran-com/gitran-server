package model

import (
	"github.com/gitran-com/gitran-server/config"
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
			tx.Create(&oldVote)
			if vote == VoteLike {
				likes_delta = 1
			} else if vote == VoteUnlike {
				unlikes_delta = 1
			}
		} else {
			if oldVote.Vote != vote {
				likes_delta = calcLikesDelta(oldVote.Vote, vote)
				unlikes_delta = calcUnlikesDelta(oldVote.Vote, vote)
				tx.Model(&oldVote).Update("vote", vote)
			}
		}
		if likes_delta != 0 || unlikes_delta != 0 {
			tx.Model(&tran).Updates(map[string]interface{}{"likes": gorm.Expr("likes+(?)", likes_delta), "unlikes": gorm.Expr("unlikes+(?)", unlikes_delta)})
		}
		likes = int(tran.Likes) + likes_delta
		return nil
	})
	return
}
