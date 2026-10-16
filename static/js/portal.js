let remedios = [];

const campoNome = document.getElementById("campoNome");
const campoRemedio = document.getElementById("campoRemedio");
const form = document.getElementById("formSenha");
const divErro = document.getElementById("erro");
const senhaGerada = document.getElementById("senhaGerada");
const busca = document.getElementById("busca");
const tabela = document.getElementById("tabelaRemedios");

async function carregarRemedios() {
  const resp = await fetch("/api/remedios");
  remedios = await resp.json();

  campoRemedio.innerHTML = "";
  remedios.forEach((r) => {
    const op = document.createElement("option");
    op.value = r.id;
    op.textContent =
      r.nome + (r.estoque <= 0 ? " (em falta)" : "");
    if (r.estoque <= 0) op.disabled = true;
    campoRemedio.appendChild(op);
  });

  mostrarTabela("");
}

function mostrarTabela(filtro) {
  filtro = filtro.toLowerCase();
  tabela.innerHTML = "";

  remedios
    .filter((r) => r.nome.toLowerCase().includes(filtro))
    .forEach((r) => {
      const linha = document.createElement("tr");
      const disponivel =
        r.estoque > 0
          ? '<span class="badge bg-success">Sim</span>'
          : '<span class="badge bg-secondary">Em falta</span>';
      linha.innerHTML =
        "<td>" + r.nome + "</td>" +
        "<td>" + r.categoria + "</td>" +
        '<td class="text-center">' + disponivel + "</td>";
      tabela.appendChild(linha);
    });
}

busca.addEventListener("input", () => mostrarTabela(busca.value));

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  divErro.classList.add("d-none");
  senhaGerada.classList.add("d-none");

  const corpo = {
    nome: campoNome.value.trim(),
    remedioId: parseInt(campoRemedio.value, 10),
  };

  const resp = await fetch("/api/senha", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(corpo),
  });
  const dados = await resp.json();

  if (!resp.ok) {
    divErro.textContent = dados.erro || "Não foi possível tirar a senha.";
    divErro.classList.remove("d-none");
    return;
  }

  document.getElementById("senhaCodigo").textContent = dados.codigo;
  document.getElementById("senhaInfo").textContent =
    dados.nome + " — " + dados.remedioNome;
  senhaGerada.classList.remove("d-none");
  form.reset();
});

carregarRemedios();
