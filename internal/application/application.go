package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MaxMindshtorm/calculator/pkg/calcul"
)

type Config struct {
	Port int
}

type Application struct {
	Cfg Config
}

func New(config Config) *Application {
	return &Application{
		Cfg: config,
	}
}

type Expression struct {
	Expression string `json:"expression"`
}

type Result struct {
	Res float64 `json:"result"`
}

type Error struct {
	Err string `json:"error"`
}

func ApplicationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Error{
			Err: "Method is not allowed",
		})
		return
	}
	var exp Expression
	err := json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Error{
			Err: "Internal server error",
		})
		return
	}

	if exp.Expression == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{
			Err: "Bad request",
		})
		return
	}

	result, err := calcul.Calc(exp.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Error{
			Err: "Expression is not valid" + ": " + fmt.Sprint(err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Result{
		Res: result,
	})
}

func (a *Application) StartServer() {
	http.HandleFunc("/api/v1/calculate", ApplicationHandler)
	http.ListenAndServe(":"+fmt.Sprint(a.Cfg.Port), nil)
}
