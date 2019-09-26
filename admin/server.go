package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	dc "github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/deputy"
	"github.com/binance-chain/bep3-deputy/util"
)

const numPerPage = 100

type Admin struct {
	Config *util.Config
	Deputy *deputy.Deputy
}

func NewAdmin(config *util.Config, deputy *deputy.Deputy) *Admin {
	return &Admin{
		Config: config,
		Deputy: deputy,
	}
}

func (admin *Admin) StatusHandler(w http.ResponseWriter, r *http.Request) {
	deputyStatus, err := admin.Deputy.Status()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.MarshalIndent(deputyStatus, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (admin *Admin) FailedSwapsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pageStr := params["page"]
	if pageStr == "" {
		http.Error(w, "required parameter 'page' is missing", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("parse page error, err=%s", err.Error()), http.StatusBadRequest)
		return
	}

	if page < 1 {
		http.Error(w, "page should be no less than 1", http.StatusBadRequest)
		return
	}

	failedSwaps, totalCount, err := admin.Deputy.FailedSwaps((page-1)*numPerPage, numPerPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := &dc.FailedSwaps{
		TotalCount: totalCount,
		CurPage:    page,
		NumPerPage: numPerPage,
		Swaps:      failedSwaps,
	}

	jsonBytes, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (admin *Admin) ResendTxHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	if idStr == "" {
		http.Error(w, "required parameter 'id' is missing", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("parse id error, err=%s", err.Error()), http.StatusBadRequest)
		return
	}

	txHash, err := admin.Deputy.ResendTx(int64(id))

	response := struct {
		TxHash string `json:"tx_hash"`
		ErrMsg string `json:"err_msg"`
	}{
		TxHash: txHash,
	}
	if err != nil {
		response.ErrMsg = err.Error()
	}

	jsonBytes, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (admin *Admin) SetModeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	modeStr := params["mode"]
	if modeStr == "" {
		http.Error(w, "required parameter 'mode' is missing", http.StatusBadRequest)
		return
	}

	mode, err := strconv.Atoi(modeStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("parse mode error, err=%s", err.Error()), http.StatusBadRequest)
		return
	}

	if dc.DeputyMode(mode) != dc.DeputyModeNormal && dc.DeputyMode(mode) != dc.DeputyModeStopSendHTLT {
		http.Error(w, fmt.Sprintf("mode only supports %d(%s) and %d(%s)",
			dc.DeputyModeNormal, dc.DeputyModeNormal, dc.DeputyModeStopSendHTLT, dc.DeputyModeStopSendHTLT), http.StatusBadRequest)
		return
	}

	admin.Deputy.SetMode(dc.DeputyMode(mode))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (admin *Admin) Endpoints(w http.ResponseWriter, r *http.Request) {
	endpoints := struct {
		Endpoints []string `json:"endpoints"`
	}{
		Endpoints: []string{
			"/status",
			"/failed_swaps/{page}",
			"/resend_tx/{id}",
			"/set_mode/{mode}",
		},
	}

	jsonBytes, err := json.MarshalIndent(endpoints, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (admin *Admin) Serve() {
	router := mux.NewRouter()

	router.HandleFunc("/", admin.Endpoints)
	router.HandleFunc("/status", admin.StatusHandler)
	router.HandleFunc("/failed_swaps/{page}", admin.FailedSwapsHandler)
	router.HandleFunc("/resend_tx/{id}", admin.ResendTxHandler)
	router.HandleFunc("/set_mode/{mode}", admin.SetModeHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         admin.Config.AdminConfig.ListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	util.Logger.Infof("start admin server at %s", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("start admin server error, err=%s", err.Error()))
	}
}
