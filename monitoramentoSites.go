package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramento = 5
const delay = 3
const arquivo = "sites.txt"
const logsRoot = "logs/"

func main() {
	leSitesArquivo()
	exibeIntroducao()
	for { //for sem parametros funciona como um while, que não existe no go

		exibeMenu()
		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			logs()
		case 0:
			fmt.Println("Saindo do programa...")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O valor da variável comando é:", comandoLido)

	return comandoLido
}

func exibeMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
}

func exibeIntroducao() {
	nome, idade := dadosIntroducao()
	versao := 1.1
	fmt.Println("Olá, sr(a).", nome, ",", idade, "anos")
	fmt.Println("Este programa está na versão", versao)
}

func dadosIntroducao() (string, int) {
	nome := "Ernani"
	idade := 40
	return nome, idade
}

func iniciarMonitoramento() {
	// var sites [4]string
	// sites[0] = "https://random-status-code.herokuapp.com/"
	// sites[1] = "https://www.alura.com.br"
	// sites[2] = "https://www.caelum.com.br"
	// fmt.Println("val", sites[3])

	// sites := []string{"https://random-status-code.herokuapp.com/", "https://www.alura.com.br"}
	// sites = append(sites, "https://www.caelum.com.br")

	sites := leSitesArquivo()

	fmt.Println("Monitorando...")
	for i := 0; i < monitoramento; i++ {
		for _, site := range sites {
			testaSite(site)
		}
		time.Sleep(delay * time.Second)
	}

}

func testaSite(site string) {
	resp, err := http.Get(site) //Caso não tenha interesse no erro, p.e., pode chamar a função assim: resp,_ := http.Get(site)
	if resp.StatusCode == 200 {
		fmt.Println("Site", site, "carregado com sucesso!")
		salvaLog(site, true)
	} else {
		fmt.Println("Erro no site", site, ". Status code:", resp.StatusCode)
		fmt.Println("Error:", err)
		salvaLog(site, false)
	}
	fmt.Println("")
}

func leSitesArquivo() []string {
	var sites []string
	arquivo, err := os.Open(arquivo)

	if err != nil {
		fmt.Println("Erro ao abrir arquivo")
	} else {
		leitor := bufio.NewReader(arquivo)

		for {
			linha, err := leitor.ReadString('\n')
			linha = strings.TrimSpace(linha)
			sites = append(sites, linha)

			if err == io.EOF {
				break
			}

		}
	}
	arquivo.Close()
	return sites
}

func salvaLog(site string, status bool) {
	year, month, day := time.Now().Date()
	logFileName := logsRoot + "log" + strconv.Itoa(year) + month.String() + strconv.Itoa(day) + ".txt"
	arquivo, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	} else {
		arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - status:" + strconv.FormatBool(status) + "\n")
	}

	arquivo.Close()

}

func listaLogs() {
	files, err := ioutil.ReadDir(logsRoot)
	fmt.Println("------ LISTA DE LOGS ------")
	if err != nil {
		fmt.Println("Erro ao imprimir arquivo")
	} else {
		for _, file := range files {
			fmt.Println(file.Name())
		}
	}
	fmt.Println("------ ------------- ------")
	fmt.Println("")

}

// restante do código omitido

func imprimeLogs(log string) {

	arquivo, err := ioutil.ReadFile(logsRoot + log)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	fmt.Println(string(arquivo))
}

func logs() {
	var log string
	listaLogs()
	fmt.Println("Digite o nome do arquivo que deseja ler (ou digite 'sair' para voltar ao menu principal):")
	fmt.Scan(&log)

	if log != "sair" {
		imprimeLogs(log)
	}
}
