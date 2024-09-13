package cmd

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var PlayerCmd = &cobra.Command{
	Use:   "player ready",
	Short: "Set the player as ready",
	Run: func(cmd *cobra.Command, args []string) {
		req, err := http.NewRequest(http.MethodPost, "http://localhost:7070/player/ready", nil)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImhleTEiLCJUb2tlbiI6bnVsbCwiZXhwIjoxNzI2MjcyNDA4fQ.5hj3viOazSYOpkHpMgJjk4Z4IILtr8htpxDMX7PdHew")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		fmt.Println("Response:", buf.String())
	},
}
