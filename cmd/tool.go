package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"

	"github.com/charlieegan3/activities-rss/pkg/tool"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	toolCfg, ok := viper.Get("tools.activities-rss").(map[string]interface{})
	if !ok {
		log.Fatalf("failed to read tools config in map[string]interface{} format")
	}
	fmt.Println(toolCfg)
	t := &tool.ActivitiesRSS{}
	t.SetConfig(toolCfg)

	j, err := t.Jobs()
	if err != nil {
		log.Fatalf("failed to get jobs: %s", err)
	}

	err = j[0].Run(context.Background())
	if err != nil {
		log.Fatalf("failed to run job: %s", err)
	}
}
