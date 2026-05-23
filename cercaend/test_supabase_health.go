package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://epawttrarbrpzmdbmxyn.supabase.co/auth/v1/health"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVwYXd0dHJhcmJycHptZGJteHluIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzY2MzA0NzMsImV4cCI6MjA5MjIwNjQ3M30.zhjusW5PGl4dRDsi38FtuFwCsfFw_wNSAG7oUuM_Dds")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Status: %d\nBody: %s\n", resp.StatusCode, string(body))
}
