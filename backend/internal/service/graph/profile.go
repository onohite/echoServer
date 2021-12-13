package graph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/labstack/gommon/log"
)

type GameProfileModel struct {
	Uid         string   `json:"uid,omitempty"`
	UserID      string   `json:"user_id,omitempty"`
	GameID      int      `json:"game_id,omitempty"`
	RankID      int      `json:"rank_id,omitempty"`
	Sex         int      `json:"sex,omitempty"`
	Age         int      `json:"age,omitempty"`
	Contact     string   `json:"contact,omitempty"`
	Description string   `json:"description,omitempty"`
	DType       []string `json:"dgraph.type,omitempty"`
}

type GameProfile struct {
	UserID      string `json:"user_id,omitempty"`
	GameID      int    `json:"game_id,omitempty"`
	RankID      int    `json:"rank_id,omitempty"`
	Sex         int    `json:"sex,omitempty"`
	Age         int    `json:"age,omitempty"`
	Contact     string `json:"contact,omitempty"`
	Description string `json:"description,omitempty"`
}

type UidRequest struct {
	Data struct {
		All []struct {
			Uid string `json:"uid"`
		} `json:"all"`
	} `json:"data"`
}

type ProfileRequest struct {
	Data struct {
		All []struct {
			UserId      string `json:"user_id"`
			GameId      int    `json:"game_id"`
			RankId      int    `json:"rank_id"`
			Age         int    `json:"age"`
			Sex         int    `json:"sex"`
			Contact     string `json:"contact"`
			Description string `json:"description"`
		} `json:"all"`
	} `json:"data"`
}

func (c GraphConn) SetProfile(p GameProfile) (string, error) {
	connCtx, cancel := context.WithTimeout(c.ctx, waitTimeout)
	defer cancel()
	txn := c.client.NewTxn()
	defer txn.Discard(connCtx)

	mutation := GameProfileModel{
		Uid:         "_:profile",
		UserID:      p.UserID,
		GameID:      p.GameID,
		RankID:      p.RankID,
		Sex:         p.Sex,
		Age:         p.Age,
		Contact:     p.Contact,
		Description: p.Description,
		DType:       []string{"Profile"},
	}

	pb, err := json.Marshal(mutation)
	if err != nil {
		log.Error(err)
		return "", err
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   pb,
	}

	res, err := txn.Mutate(connCtx, mu)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Printf("SetProfile complete successful with %v", res)

	query := fmt.Sprintf(`query {
    all(func: eq(user_id, %s)) {
      uid
    }
}`, p.UserID)

	secondTxn := c.client.NewTxn()
	defer secondTxn.Discard(connCtx)

	resp, err := secondTxn.Query(connCtx, query)
	if err != nil {
		log.Error(err)
		return "", err
	}
	var uidReq UidRequest
	log.Printf("%v", resp)
	err = json.Unmarshal(resp.GetJson(), &uidReq)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Printf("found uid %s", uidReq.Data.All[0].Uid)
	return uidReq.Data.All[0].Uid, nil
}

func (c GraphConn) GetProfile(uid string) (*GameProfile, error) {
	connCtx, cancel := context.WithTimeout(c.ctx, waitTimeout)
	defer cancel()
	txn := c.client.NewTxn()
	defer txn.Discard(connCtx)

	query := `query all($a: string){
    all(func: uid($a)) {
		user_id
    	game_id
    	rank_id
    	age
    	sex
    	contact
    	description
    }
}`

	resp, err := txn.QueryWithVars(connCtx, query, map[string]string{"$a": uid})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var profile ProfileRequest
	err = json.Unmarshal(resp.Json, &profile)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Printf("find complete successful with %v", profile)
	copyData := profile.Data.All[0]
	var respProfile GameProfile
	if len(profile.Data.All) < 1 {
		return nil, errors.New("по данному uid профиль не найден")
	}
	respProfile.UserID = copyData.UserId
	respProfile.GameID = copyData.GameId
	respProfile.RankID = copyData.RankId
	respProfile.Age = copyData.Age
	respProfile.Sex = copyData.Sex
	respProfile.Contact = copyData.Contact
	respProfile.Description = copyData.Description
	return &respProfile, nil
}
