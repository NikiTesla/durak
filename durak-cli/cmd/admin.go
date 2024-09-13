package cmd

import (
	"bytes"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin start_game",
	Short: "Start the game as admin",
	Run: func(cmd *cobra.Command, args []string) {
		req, err := http.NewRequestWithContext(cmd.Context(), http.MethodPost, "http://localhost:7070/admin/start_game", nil)
		if err != nil {
			log.WithError(err).Error("creating request")
			return
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiVG9rZW4iOm51bGwsImV4cCI6MTcyNjI3MjY3Nn0.P5B1cqSmL3uV9kon2um8lJXTCr8h9dwsRKNfeTYn1nY")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.WithError(err).Error("executing request")
			return
		}
		defer resp.Body.Close()

		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		fmt.Println("Response:", buf.String())
	},
}
