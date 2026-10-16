package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
)

const porta = "3000"
const arquivo = "dados.json"

type Remedio struct {
	ID            int    `json:"id"`
	Nome          string `json:"nome"`
	Categoria     string `json:"categoria"`
	Estoque       int    `json:"estoque"`
	EstoqueMinimo int    `json:"estoqueMinimo"`
}

type Senha struct {
	Codigo      string `json:"codigo"`
	Nome        string `json:"nome"`
	RemedioID   int    `json:"remedioId"`
	RemedioNome string `json:"remedioNome"`
	Status      string `json:"status"` // "aguardando" ou "chamada"
	Guiche      int    `json:"guiche"`
}

type Dados struct {
	Remedios   []Remedio `json:"remedios"`
	Senhas     []Senha   `json:"senhas"`
	ProximoNum int       `json:"proximoNum"`
}

var (
	mu    sync.Mutex
	dados Dados
)

// carrega os dados do arquivo JSON; se nao existir, usa a lista inicial.
func carregar() {
	conteudo, err := os.ReadFile(arquivo)
	if err != nil {
		dados = dadosIniciais()
		salvar()
		return
	}
	if err := json.Unmarshal(conteudo, &dados); err != nil {
		dados = dadosIniciais()
		salvar()
	}
}

// salva os dados atuais no arquivo JSON.
func salvar() {
	conteudo, _ := json.MarshalIndent(dados, "", "  ")
	os.WriteFile(arquivo, conteudo, 0644)
}

// lista inicial de remedios da Farmacia Popular.
func dadosIniciais() Dados {
	remedios := []Remedio{
		{1, "Losartana Potássica 50mg", "Hipertensão", 120, 30},
		{2, "Captopril 25mg", "Hipertensão", 25, 30},
		{3, "Enalapril 10mg", "Hipertensão", 80, 30},
		{4, "Hidroclorotiazida 25mg", "Hipertensão", 60, 20},
		{5, "Atenolol 25mg", "Hipertensão", 15, 20},
		{6, "Propranolol 40mg", "Hipertensão", 90, 20},
		{7, "Cloridrato de Metformina 500mg", "Diabetes", 200, 50},
		{8, "Glibenclamida 5mg", "Diabetes", 40, 40},
		{9, "Insulina NPH", "Diabetes", 18, 25},
		{10, "Insulina Regular", "Diabetes", 22, 25},
		{11, "Sulfato de Salbutamol (spray)", "Asma", 35, 15},
		{12, "Dipropionato de Beclometasona", "Asma", 10, 15},
		{13, "Brometo de Ipratrópio", "Asma", 28, 15},
		{14, "Sinvastatina 20mg", "Colesterol", 70, 30},
		{15, "Ácido Acetilsalicílico 100mg", "Prevenção", 150, 40},
		{16, "Omeprazol 20mg", "Gástrico", 55, 30},
	}
	return Dados{Remedios: remedios, Senhas: []Senha{}, ProximoNum: 1}
}

// acha um remedio pelo id (retorna o indice no slice, ou -1).
func indiceRemedio(id int) int {
	for i, r := range dados.Remedios {
		if r.ID == id {
			return i
		}
	}
	return -1
}

// ----- Handlers (as rotas da API) -----

// GET /api/remedios -> lista todos os remedios (ordenados por nome).
func handlerRemedios(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	lista := make([]Remedio, len(dados.Remedios))
	copy(lista, dados.Remedios)
	sort.Slice(lista, func(i, j int) bool { return lista[i].Nome < lista[j].Nome })

	responderJSON(w, lista)
}

// POST /api/senha -> cria uma senha para uma pessoa pegar um remedio.
// Corpo esperado: {"nome": "...", "remedioId": 1}
func handlerCriarSenha(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var corpo struct {
		Nome      string `json:"nome"`
		RemedioID int    `json:"remedioId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&corpo); err != nil {
		responderErro(w, "Dados inválidos.", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if corpo.Nome == "" {
		responderErro(w, "Digite seu nome.", http.StatusBadRequest)
		return
	}

	i := indiceRemedio(corpo.RemedioID)
	if i == -1 {
		responderErro(w, "Remédio não encontrado.", http.StatusNotFound)
		return
	}
	if dados.Remedios[i].Estoque <= 0 {
		responderErro(w, "Esse remédio está em falta no momento.", http.StatusConflict)
		return
	}

	senha := Senha{
		Codigo:      fmt.Sprintf("A%03d", dados.ProximoNum),
		Nome:        corpo.Nome,
		RemedioID:   dados.Remedios[i].ID,
		RemedioNome: dados.Remedios[i].Nome,
		Status:      "aguardando",
	}
	dados.ProximoNum++
	dados.Senhas = append(dados.Senhas, senha)
	salvar()

	responderJSON(w, senha)
}

// GET /api/fila -> senhas que ainda estao aguardando.
func handlerFila(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	fila := []Senha{}
	for _, s := range dados.Senhas {
		if s.Status == "aguardando" {
			fila = append(fila, s)
		}
	}
	responderJSON(w, fila)
}

// POST /api/chamar -> chama a proxima senha da fila (entrega o remedio).
// Baixa o estoque e avisa se o remedio ficou acabando.
func handlerChamar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// acha a primeira senha aguardando.
	idx := -1
	for i, s := range dados.Senhas {
		if s.Status == "aguardando" {
			idx = i
			break
		}
	}
	if idx == -1 {
		responderJSON(w, map[string]any{"vazio": true})
		return
	}

	// guiche vai de 1 a 4, alternando.
	guiche := (dados.ProximoNum % 4) + 1
	dados.Senhas[idx].Status = "chamada"
	dados.Senhas[idx].Guiche = guiche

	// baixa o estoque do remedio.
	alerta := false
	estoqueAtual := 0
	ri := indiceRemedio(dados.Senhas[idx].RemedioID)
	if ri != -1 {
		if dados.Remedios[ri].Estoque > 0 {
			dados.Remedios[ri].Estoque--
		}
		estoqueAtual = dados.Remedios[ri].Estoque
		if estoqueAtual <= dados.Remedios[ri].EstoqueMinimo {
			alerta = true
		}
	}

	salvar()

	responderJSON(w, map[string]any{
		"senha":        dados.Senhas[idx],
		"alerta":       alerta,
		"estoqueAtual": estoqueAtual,
	})
}

// GET /api/atual -> a ultima senha que foi chamada (pro painel mostrar).
func handlerAtual(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var atual *Senha
	for i := len(dados.Senhas) - 1; i >= 0; i-- {
		if dados.Senhas[i].Status == "chamada" {
			atual = &dados.Senhas[i]
			break
		}
	}

	if atual == nil {
		responderJSON(w, map[string]any{"vazio": true})
		return
	}

	// verifica se o remedio dessa senha esta acabando.
	alerta := false
	estoqueAtual := 0
	ri := indiceRemedio(atual.RemedioID)
	if ri != -1 {
		estoqueAtual = dados.Remedios[ri].Estoque
		if estoqueAtual <= dados.Remedios[ri].EstoqueMinimo {
			alerta = true
		}
	}

	responderJSON(w, map[string]any{
		"senha":        atual,
		"alerta":       alerta,
		"estoqueAtual": estoqueAtual,
	})
}

// GET /api/estoque -> remedios + se o estoque esta baixo (pro dashboard).
func handlerEstoque(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	type Item struct {
		Remedio
		Baixo bool `json:"baixo"`
	}

	itens := []Item{}
	total, baixos := 0, 0
	for _, rem := range dados.Remedios {
		baixo := rem.Estoque <= rem.EstoqueMinimo
		if baixo {
			baixos++
		}
		total += rem.Estoque
		itens = append(itens, Item{rem, baixo})
	}

	sort.Slice(itens, func(i, j int) bool { return itens[i].Nome < itens[j].Nome })

	responderJSON(w, map[string]any{
		"itens":        itens,
		"totalUnidades": total,
		"qtdRemedios":  len(dados.Remedios),
		"qtdBaixos":    baixos,
	})
}

// ----- Ajudantes -----

func responderJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func responderErro(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"erro": msg})
}

// serve uma pagina HTML da pasta static.
func servirPagina(nome string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/"+nome)
	}
}

func main() {
	carregar()

	mux := http.NewServeMux()

	// paginas
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "static/index.html")
	})
	mux.HandleFunc("/painel", servirPagina("painel.html"))
	mux.HandleFunc("/dashboard", servirPagina("dashboard.html"))

	// arquivos estaticos (css e js)
	mux.Handle("/css/", http.FileServer(http.Dir("static")))
	mux.Handle("/js/", http.FileServer(http.Dir("static")))

	// API
	mux.HandleFunc("/api/remedios", handlerRemedios)
	mux.HandleFunc("/api/senha", handlerCriarSenha)
	mux.HandleFunc("/api/fila", handlerFila)
	mux.HandleFunc("/api/chamar", handlerChamar)
	mux.HandleFunc("/api/atual", handlerAtual)
	mux.HandleFunc("/api/estoque", handlerEstoque)

	fmt.Printf("Servidor rodando! Abra no navegador: http://localhost:%s\n", porta)
	log.Fatal(http.ListenAndServe(":"+porta, mux))
}
