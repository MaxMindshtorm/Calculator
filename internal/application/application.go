package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/MaxMindshtorm/calculator/pkg/calcul"
)

type Config struct {
	Port string
}

func ConfigFromEnv() *Config {
	cfg := new(Config)
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	return cfg
}

type Application struct {
	Cfg Config
}

func New() *Application {
	return &Application{
		Cfg: *ConfigFromEnv(),
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
	err := http.ListenAndServe("localhost:"+fmt.Sprint(a.Cfg.Port), nil)
	if err != nil {
		errors.New("Error while starting srever")
	}
}
