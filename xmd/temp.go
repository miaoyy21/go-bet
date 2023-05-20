package xmd

import (
	"encoding/json"
	"os"
)

type TempData struct {
	Issue    int
	Latest   map[int]struct{}
	Rts      map[int]float64
	Surplus  int
	BetGold  int
	Exp      float64
	Dev      float64
	UserGold int
}

func tempSave() error {
	data := &TempData{
		Issue:    issue,
		Latest:   latest,
		Rts:      xRts,
		Surplus:  xSurplus,
		BetGold:  xBetGold,
		Exp:      xExp,
		Dev:      xDev,
		UserGold: xUserGold,
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

	issue = data.Issue
	latest = data.Latest
	xRts = data.Rts
	xSurplus = data.Surplus
	xBetGold = data.BetGold
	xExp = data.Exp
	xDev = data.Dev
	xUserGold = data.UserGold

	return nil
}
