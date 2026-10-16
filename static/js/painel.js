const btnChamar = document.getElementById("btnChamar");
const alertaEstoque = document.getElementById("alertaEstoque");

function mostrarSenha(info) {
  if (!info || info.vazio || !info.senha) {
    document.getElementById("senhaCodigo").textContent = "---";
    document.getElementById("senhaNome").textContent = "";
    document.getElementById("senhaRemedio").textContent = "";
    document.getElementById("senhaGuiche").classList.add("d-none");
    return;
  }

  const s = info.senha;
  document.getElementById("senhaCodigo").textContent = s.codigo;
  document.getElementById("senhaNome").textContent = s.nome;
  document.getElementById("senhaRemedio").textContent = s.remedioNome;

  const guiche = document.getElementById("senhaGuiche");
  guiche.textContent = "Guichê " + s.guiche;
  guiche.classList.remove("d-none");

  // Aviso de estoque acabando para o remédio que acabou de ser chamado.
  if (info.alerta) {
    alertaEstoque.textContent =
      "⚠️ Estoque de " + s.remedioNome + " está acabando (" +
      info.estoqueAtual + " restantes). Reabasteça!";
    alertaEstoque.classList.remove("d-none");
  } else {
    alertaEstoque.classList.add("d-none");
  }
}

async function atualizarFila() {
  const resp = await fetch("/api/fila");
  const fila = await resp.json();
  document.getElementById("qtdFila").textContent = fila.length;
}

async function atualizarAtual() {
  const resp = await fetch("/api/atual");
  const info = await resp.json();
  mostrarSenha(info);
}

btnChamar.addEventListener("click", async () => {
  const resp = await fetch("/api/chamar", { method: "POST" });
  const info = await resp.json();

  if (info.vazio) {
    alertaEstoque.textContent = "Não há ninguém na fila.";
    alertaEstoque.className = "alert alert-info text-center";
    setTimeout(() => alertaEstoque.classList.add("d-none"), 2500);
    return;
  }

  alertaEstoque.className = "alert alert-warning text-center";
  mostrarSenha(info);
  atualizarFila();
});

// Atualiza automaticamente a cada 3 segundos (painel "ao vivo").
function cicloAutomatico() {
  atualizarAtual();
  atualizarFila();
}
cicloAutomatico();
setInterval(cicloAutomatico, 3000);
