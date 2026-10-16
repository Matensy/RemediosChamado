# 💊 Remédios Chamados — Farmácia Popular

Sistema para uma farmácia popular com três partes conectadas:

- **Portal** — a pessoa consulta a lista de remédios disponíveis, escolhe um
  e **tira uma senha** (com o nome e o remédio).
- **Painel de Senhas** — chama a senha mostrando **nome + remédio + guichê**,
  e **avisa quando aquele remédio está acabando**.
- **Dashboard** — mostra o **estoque** de cada remédio, destacando os que
  estão no fim.

Quando uma senha é chamada (remédio entregue), o estoque daquele remédio
diminui automaticamente.

## 🛠️ Tecnologias usadas

- **Back-end:** [Go (Golang)](https://go.dev/) — usando só a biblioteca padrão
- **Front-end:** HTML + [Bootstrap 5](https://getbootstrap.com/) + JavaScript
- **Dados:** salvos em um arquivo `dados.json` (criado automaticamente)

## 📁 Estrutura dos arquivos

```
RemdiosChamados/
├─ main.go              → o back-end (servidor + API + regras de estoque/senhas)
├─ go.mod               → identifica o módulo Go
└─ static/              → o site
   ├─ index.html        → Portal (tirar senha)
   ├─ painel.html       → Painel de senhas
   ├─ dashboard.html    → Dashboard de estoque
   ├─ css/style.css
   └─ js/
      ├─ portal.js
      ├─ painel.js
      └─ dashboard.js
```

## ▶️ Como rodar o projeto

> Antes de começar, você precisa ter o [Go](https://go.dev/dl/) instalado.

1. Abra o terminal dentro da pasta do projeto.

2. Rode o servidor:

   ```bash
   go run .
   ```

3. Abra no navegador:

   - Portal: `http://localhost:3000`
   - Painel: `http://localhost:3000/painel`
   - Dashboard: `http://localhost:3000/dashboard`

Na primeira vez, o arquivo `dados.json` é criado sozinho com a lista inicial
de remédios.

## 🧭 Como funciona (passo a passo)

1. No **Portal**, a pessoa digita o nome, escolhe o remédio e tira a senha.
2. A senha entra na **fila** (status "aguardando").
3. No **Painel**, o atendente clica em **Chamar próximo**: a senha aparece
   em destaque (nome, remédio e guichê) e o **estoque daquele remédio baixa 1**.
4. Se o estoque ficar no limite mínimo, o painel mostra um **aviso de que o
   remédio está acabando**.
5. O **Dashboard** mostra, em tempo quase real, quanto ainda tem de cada
   remédio e quais estão acabando.

## 💡 Observações

- Os dados ficam em `dados.json`. Se quiser **zerar** o estoque e as senhas,
  é só apagar esse arquivo e rodar de novo.
- A lista de remédios inicial está no próprio `main.go` (função
  `dadosIniciais`), fácil de editar.

---

Feito para estudo. 🚀
