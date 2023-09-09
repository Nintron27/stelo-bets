package routes

import (
	"bytes"
	"chi-learning/internal/env"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func InternalRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", postTransaction)

	return r
}

type transaction struct {
	WalletId int64  `json:"wallet_id" validate:"required"`
	Memo     string `json:"memo"`
	Assets   asset  `json:"assets" validate:"required"`
}

type asset struct {
	Stelo uint64 `json:"stelo" validate:"required,min=1000,max=1000000"`
}

func postTransaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("key") != env.Variables.ApiKey {
		http.Error(w, "Invalid API key", http.StatusForbidden)
		return
	}

	var tx transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chanceInput, err := strconv.Atoi(tx.Memo)
	if err != nil || !(chanceInput >= 2 && chanceInput <= 100) {
		http.Error(w, "Invalid memo", http.StatusBadRequest)
		return
	}

	edgedChance := 1.0/float64(chanceInput) - 0.05/float64(chanceInput)
	roll := rand.Float64()

	if roll < edgedChance {
		w.WriteHeader(http.StatusAccepted)
		go sendPayout(fmt.Sprint(tx.WalletId), chanceInput, uint64(chanceInput)*tx.Assets.Stelo)
	}

	w.WriteHeader(http.StatusOK)
}

// Send payout to Stelo Finance API.
// Should be run in a go-routine so that way
// the transactions aren't blocked by the DB lock
func sendPayout(walletId string, multiplier int, payout uint64) {
	body, err := json.Marshal(map[string]interface{}{
		"recipient": walletId,
		"type":      2,
		"memo":      fmt.Sprintf("You won %vx, congrats!", multiplier),
		"assets": map[string]interface{}{
			"stelo": payout,
		},
	})
	if err != nil {
		fmt.Printf("Failed to marshal json, payout (%v), to: %v\n", payout, walletId)
		return
	}

	req, err := http.NewRequest(http.MethodPost, env.Variables.SteloApi+"/wallet/transactions", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error creating req, payout (%v), to: %v\n", payout, walletId)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.Variables.WalletKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending payout (%v), to: %v\n", payout, walletId)
		return
	}

	resp.Body.Close()
}
