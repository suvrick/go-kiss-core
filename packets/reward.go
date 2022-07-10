package packets

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var rewards []string

type Caption struct {
	Ru string `json:"ru"`
	En string `json:"en"`
}

type Obj struct {
	ID       uint64  `json:"id"`
	Captions Caption `json:"captions"`
}

func GetReward(rewardID uint64) string {
	pattern := fmt.Sprintf("\"id\": %d,", rewardID)
	for _, r := range rewards {
		if strings.Index(r, pattern) > 0 {
			return r
		}
	}

	return ""
}

func init_rewards() error {

	resp, err := http.Get("https://bottleconf.realcdn.ru/rewards.json")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("reward bad request")
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	jsons := strings.ReplaceAll(string(body), "\t", "")
	lines := strings.Split(jsons, "\n")

	rewards = make([]string, 0)

	for _, line := range lines {
		if strings.Index(line, "id") > 0 && strings.Index(line, "captions") > 0 && strings.Index(line, "content") > 0 {

			if line[len(line)-1] == ',' {
				line = line[:len(line)-1]
			}

			rewards = append(rewards, line)
		}
	}

	return nil
}
