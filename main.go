package main

//Importando a biblioteca
import (
	"encoding/json" // permite trabalhar com JSON (formato de troca de dados, parecido com texto).
	"fmt"           // para imprimir mensagens no terminal, como “servidor rodando…”.
	"net/http"      //para criar o servidor e definir as rotas da API.
)

// Define a estrutura dos dados que vamos receber
type OperacaoRequest struct {
	Operando1 float64 `json:"operando1"`
	Operando2 float64 `json:"operando2"`
}

// Define o formato da resposta
type ResultadoResponse struct {
	Resultado float64 `json:"resultado"`
}

// Cada linha dessas define uma rota da API
func main() {
	http.HandleFunc("/soma", somaHandler)
	http.HandleFunc("/subtracao", subtracaoHandler)
	http.HandleFunc("/multiplicacao", multiplicacaoHandler)
	http.HandleFunc("/divisao", divisaoHandler)
	//fmt.Println(...): mostra no terminal que o servidor começou
	//http.ListenAndServe(...): inicia o servidor na porta 8080 e espera receber requisições
	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}

// Quando a rota /soma for acessada, essa função:
// Usa uma função chamada executarOperacao
// E diz o que ela faz com os dois números: a + b (soma).
func somaHandler(w http.ResponseWriter, r *http.Request) {
	executarOperacao(w, r, func(a, b float64) float64 {
		return a + b
	})
}

func subtracaoHandler(w http.ResponseWriter, r *http.Request) {
	executarOperacao(w, r, func(a, b float64) float64 {
		return a - b
	})
}

func multiplicacaoHandler(w http.ResponseWriter, r *http.Request) {
	executarOperacao(w, r, func(a, b float64) float64 {
		return a * b
	})
}

func divisaoHandler(w http.ResponseWriter, r *http.Request) {
	executarOperacao(w, r, func(a, b float64) float64 {
		if b == 0 {
			http.Error(w, "Erro: Divisão por zero", http.StatusBadRequest)
			return 0
		}
		return a / b
	})
}

// Funcao que executa a operacao
// w: para responder à pessoa que fez o pedido
// r: a requisição que chegou
// operacao: uma função que diz o que fazer com os dois números (somar, subtrair, etc)
// Verifica se a requisição é do tipo POST (o tipo certo para enviar dados). Se não for, devolve erro "Método não permitido"
func executarOperacao(w http.ResponseWriter, r *http.Request, operacao func(float64, float64) float64) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	//Cria a variável req para guardar os dados recebidos
	//Usa json.NewDecoder().Decode() para ler os dados JSON que vieram do corpo da requisição e colocar dentro da estrutura OperacaoRequest
	var req OperacaoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	//Se tiver algum erro ao ler os dados, retorna erro
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	// Se for divisão e operando2 for zero, a verificação está dentro do handler específico.
	//Executa a operação matemática, seja ela soma, subtração, etc., com os dois números enviados.
	resultado := operacao(req.Operando1, req.Operando2)

	//Cria a estrutura com o resultado
	//Define que a resposta será em formato JSON
	//Envia o resultado como resposta para quem fez a requisição
	resposta := ResultadoResponse{Resultado: resultado}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resposta)
}
