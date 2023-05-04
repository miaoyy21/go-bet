package xmd

import (
	"encoding/json"
	"os"
)

type TempData struct {
	Latest   map[int]struct{}
	Surplus  int
	BetGold  int
	Rx       float64
	Dx       float64
	UserGold int

	Wins  int
	Fails int
}

func tempSave() error {
	data := &TempData{
		Latest:   latest,
		Surplus:  xSurplus,
		BetGold:  xBetGold,
		Rx:       xRx,
		Dx:       xDx,
		UserGold: xUserGold,

		Wins:  wins,
		Fails: fails,
	}

	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile("temp_data.json", bs, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func tempLoad() error {
	var data TempData

	bs, err := os.ReadFile("temp_data.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bs, &data); err != nil {
		return err
	}

	latest = data.Latest
	xSurplus = data.Surplus
	xBetGold = data.BetGold
	xRx = data.Rx
	xDx = data.Dx
	xUserGold = data.UserGold
	wins = data.Wins
	fails = data.Fails

	return nil
}
