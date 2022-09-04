package frame

import (
	"encoding/json"
	"net/http"

	"github.com/suvrick/go-kiss-core/frame"
	"github.com/suvrick/go-kiss-core/game"
	"github.com/suvrick/go-kiss-core/packets/client"
)

type FrameController struct {
	h *http.ServeMux
}

func NewFrameController(mux *http.ServeMux) *FrameController {

	s := &FrameController{
		h: mux,
	}

	s.h.HandleFunc("/frame", frameHandle)
	s.h.HandleFunc("/frame/parse", parseHandle)

	return s
}

func parseHandle(w http.ResponseWriter, r *http.Request) {

	type f struct {
		Frame string `json:"frame"`
	}

	fr := f{}

	json.NewDecoder(r.Body).Decode(&fr)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	fra := frame.NewDefaultFrame()
	d, _ := fra.Parse(fr.Frame)

	g := game.NewGame(game.GetDefaultGameConfig())
	g.Run()

	login := &client.Login{
		ID:          d["login_id"].(uint64),
		NetType:     d["frame_type"].(uint16),
		DeviceType:  6,
		Key:         d["token"].(string),
		OAuth:       1,
		AccessToken: d["token2"].(string),
		Gender:      0,
	}

	g.Send(client.LOGIN, login)
	bot := <-g.Done
	json.NewEncoder(w).Encode(&bot)
}

func frameHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/static/frame/frame.html")
}
