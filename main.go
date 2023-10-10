package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/tools/go/analysis/passes/defers"
)


type Server struct {
	ServerNmae 	string
	ServerURL 	string
	tempoExecucao 	float64
	status int
	dataFalha string
}

func criarListaServidores (serverList *os.File) []Server{
	csvReader := csv.NewReader(serverList)
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var servidores []Server
	for i, line := range data {
		if i > 0 {
			servidor := Server{
				ServerNmae: line[0],
				ServerURL: line[1],
			}
			servidores = append(servidores, servidor)
		}
	}
	return servidores
}

func checkServer (servidores []Server []Server) {
	var downServers []Server
	now := timne.Now()
	for _, servidor := range servidores {
		agora := time.Now()
		get, err := http.Get(servidor.ServerURL)
		if err != nil {
			fmt. Printf("Server %s is down [$s]\n, "servidor.ServerName, err.Error())
			servidor.status = 0
			servidor.dataFalha = agora.Format("02/01/2006 15:04:05")
			downServers = append(downServers, servidor)
			continue
		}
		status := get.StatusCode
		if servidor.status != 200 {
			servidor.dataFalha = agora.Format("02/01/2006 15:04:05")
			downServers = append(downServers, servidor)
		}
		servidor.tempoExecucao = time.Since(agora).Seconds()
		fmt.Printf("Status: [%d] Tempo de carga: [%f] URL: [%s]\n", status, servidor.tempoExecucao, servidor.ServerURL)
	}
}

func openFiles(serverListFile string, downtimeFile, string)(*os.File, *os.File){
	serverList, err := os.OpenFile(serverListFile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	downtimeList, err := os.OpenFile(downtimeFile, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return serverList, downtimeList
}

func generateDowntime(downtime *os.File, downServers []Server){
	csvWriter := csv.NewWriter(downtimeList)
	for _, servidor := range downServers {
		line := []string{servidor.ServerName, servidor.ServerURL, servidor.dataFalha, fmt.Sprintf("%f", servidor.tempoExecucao), fmt.Sprintf("%d", servidor.status)}
		csvWriter.Write(line)
	}
	csvWriter.Flush()
}

func main(){
	serverList, downtimeList := openFiles(os.Args[1], os.Args[2])
	defer serverList.Close()
	defer downtimeList.Close()
	servidores := criarListaServidores(serverList)
	downServers := checkServer(servidores)
	generateDowntime(downtimeList, downServers)
}