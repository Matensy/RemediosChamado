const lista = document.getElementById("listaEstoque");

async function carregar() {
  const resp = await fetch("/api/estoque");
  const dados = await resp.json();

  document.getElementById("qtdRemedios").textContent = dados.qtdRemedios;
  document.getElementById("totalUnidades").textContent = dados.totalUnidades;
  document.getElementById("qtdBaixos").textContent = dados.qtdBaixos;

  // maior estoque serve de referência para o tamanho das barras.
  const maxEstoque = Math.max(1, ...dados.itens.map((i) => i.estoque));

  lista.innerHTML = "";
  dados.itens.forEach((item) => {
    const largura = Math.round((item.estoque / maxEstoque) * 100);
    const cor = item.baixo ? "bg-danger" : "bg-success";
    const aviso = item.baixo
      ? '<span class="badge bg-danger ms-2">Acabando</span>'
      : "";

    const bloco = document.createElement("div");
    bloco.className = "mb-3";
    bloco.innerHTML =
      '<div class="d-flex justify-content-between">' +
      '<span><strong>' + item.nome + "</strong> " +
      '<span class="text-muted small">(' + item.categoria + ")</span>" + aviso + "</span>" +
      '<span class="text-muted">' + item.estoque + " un.</span>" +
      "</div>" +
      '<div class="progress" style="height: 12px;">' +
      '<div class="progress-bar ' + cor + '" style="width: ' + largura + '%"></div>' +
      "</div>";
    lista.appendChild(bloco);
  });
}

carregar();
// Atualiza a cada 5 segundos para refletir as entregas feitas no painel.
setInterval(carregar, 5000);
