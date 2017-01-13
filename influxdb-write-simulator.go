package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client"
)

func main() {
	interval, err := time.ParseDuration(os.Getenv("WRITE_INTERVAL"))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	influxClient := GetClient()

	for {
		WriteRandomData(influxClient)
		fmt.Println("wrote points")
		time.Sleep(interval)
	}
}

func GetClient() *client.Client {
	influxdbHost := os.Getenv("INFLUXDB_HOST")
	influxdbPort, err := strconv.Atoi(os.Getenv("INFLUXDB_PORT"))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	host, err := url.Parse(fmt.Sprintf("http://%s:%d", influxdbHost, influxdbPort))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *host,
		Username: os.Getenv("INFLUXDB_USER"),
		Password: os.Getenv("INFLUXDB_PASSWD"),
	}

	influxClient, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}
	return influxClient
}

func WriteRandomData(influxClient *client.Client) {
	tables := []string{"table1", "table2"}
	columns := []string{"col1", "col2"}
	transactions := []string{"INSERT", "UPDATE", "DELETE"}
	pts := make([]client.Point, 0)
	rand.Seed(time.Now().UTC().UnixNano())

	for _, table := range tables {
		for _, col := range columns {
			pts = append(pts, client.Point{
				Measurement: "transactions",
				Tags: map[string]string{
					"table":       table,
					"column":      col,
					"transaction": transactions[rand.Intn(len(transactions))],
				},
				Fields: map[string]interface{}{
					"count": rand.Intn(10),
				},
				Precision: "s",
			})
		}

	}

	bps := client.BatchPoints{
		Points:          pts,
		Database:        "oracle-metrics",
		RetentionPolicy: "default",
	}
	_, err := influxClient.Write(bps)
	if err != nil {
		log.Fatal(err)
	}
}
