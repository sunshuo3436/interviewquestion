//"http://localhost:8080/on?id=1"查询id为1的灯泡是否打开
// "http://localhost:8080/off?id=1"查询灯泡是否关闭
// "http://localhost:8080/status?id=1"查询灯泡状态

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// 灯泡的状态
type Bulb struct {
	ID     int  `json:"id"`
	Status bool `json:"status"` // ture打开，false关上
}

var bulbs []Bulb

func turnOn(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > 80 {
		http.Error(w, "Invalid bulb ID", http.StatusBadRequest)
		return
	}

	bulbs[id-1].Status = true

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bulb %d is now ON", id)
}

func turnOff(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > 80 {
		http.Error(w, "Invalid bulb ID", http.StatusBadRequest)
		return
	}

	bulbs[id-1].Status = false

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bulb %d is now OFF", id)
}

func status(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > 80 {
		http.Error(w, "Invalid bulb ID", http.StatusBadRequest)
		return
	}

	bulb := bulbs[id-1]

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bulb)
}

func main() {
	//初始化
	bulbs = make([]Bulb, 80)
	for i := range bulbs {
		bulbs[i] = Bulb{ID: i + 1, Status: false}
	}

	http.HandleFunc("/on", turnOn)
	http.HandleFunc("/off", turnOff)
	http.HandleFunc("/status", status)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
