---
inclusion: always
---

# Documento de Steering – Produto **Ranking**

## 1. Visão Geral do Produto

**Ranking** é um site que permite que pessoas criem listas comparativas (*rankings*) e atribuam notas a itens dentro desses rankings.
Os usuários podem:

* Criar seus próprios rankings;
* Adicionar itens a rankings;
* Avaliar itens em rankings criados por si ou por outras pessoas;
* Visualizar rankings públicos ou privados;
* Compartilhar rankings com outras pessoas.

O objetivo do produto é ser **simples, rápido e sem barreiras de entrada**: não existe cadastro nem login tradicional. O sistema identifica cada usuário automaticamente no próprio navegador.

---

## 2. Conceito de Ranking

Um **ranking** é uma estrutura que organiza itens que serão avaliados com base em critérios definidos pelo criador.

Cada ranking é composto por:

* **Nome** – obrigatório
* **Descrição** – opcional
* **Visibilidade** – pública ou privada
* **Tags** – palavras-chave opcionais para facilitar busca
* **Atributos** – critérios de avaliação (máximo de 6)
* **Itens** – elementos que serão avaliados

### Exemplo de Ranking

* **Nome:** Vídeo Games
* **Descrição:** Ranking para avaliação de vídeo games
* **Tags:** “vídeo games”, “games”
* **Atributos:**

  * Gráfico
  * Som
  * Controle
  * Biblioteca de Jogos
  * Preço

Neste exemplo, cada videogame será avaliado com base nesses cinco atributos.

---

## 3. Atributos

Os **atributos** representam os critérios de avaliação.

Cada atributo possui:

* **Nome** – obrigatório
* **Descrição** – opcional

Regras:

* Um ranking pode ter **no máximo 6 atributos**.
* Apenas o **criador do ranking** pode definir e alterar os atributos.
* Os atributos são iguais para todos os usuários que avaliam aquele ranking.

Exemplo:

| Nome do Atributo    | Descrição                       |
| ------------------- | ------------------------------- |
| Gráfico             | Qualidade visual geral          |
| Biblioteca de Jogos | Quantidade e variedade de jogos |

---

## 4. Itens

Os **itens** são os elementos que serão avaliados dentro de um ranking.

Cada item possui:

* **Nome** – obrigatório
* **Foto** – opcional

Regras:

* Tanto o criador do ranking quanto qualquer usuário com acesso ao ranking podem adicionar itens.
* Se um item não tiver foto:

  * Deve ser exibida uma imagem padrão;
  * Essa imagem padrão deve conter as iniciais do nome do item.
  * Exemplo: “Sega Mega Drive” → imagem com “SM”.

---

## 5. Visualização do Ranking

Quando um usuário acessa um ranking, ele o visualiza em formato de **tabela**:

* Cada **linha** representa um item;
* Cada **coluna** representa um atributo;
* Existe uma coluna final chamada **Overall**.

Exemplo:

| **Vídeo Games**                   | Gráfico | Som | Controle | Biblioteca de Jogos | Preço | Overall |
| --------------------------------- | ------- | --- | -------- | ------------------- | ----- | ------- |
| **Sega Mega Drive**               | 79      | 83  | 80       | 92                  | 89    | 84,6    |
| **Nintendo Entertainment System** | 88      | 80  | 87       | 94                  | 90    | 87,8    |
| **Neo Geo CD**                    | 99      | 97  | 90       | 65                  | 49    | 80,0    |

Regras de exibição:

* Na primeira coluna:

  * Deve aparecer a foto do item (se existir);
  * Abaixo da foto, o nome do item.
* A coluna **Overall**:

  * É calculada automaticamente;
  * Representa a **média aritmética** das notas dos atributos;
  * Não é armazenada permanentemente;
  * Sempre é recalculada no momento da visualização.

---

## 6. Modos de Visualização

O usuário pode alternar entre dois modos:

1. **Média de notas de todos os usuários**

   * Mostra, para cada item e atributo, a média das notas dadas por todas as pessoas.
   * O overall também é uma média baseada nessas notas agregadas.

2. **Minhas notas**

   * Mostra apenas as notas atribuídas pelo próprio usuário.
   * O usuário pode inserir valores de **0 a 100** em cada atributo de cada item.

Regras no modo “Minhas notas”:

* O usuário só pode dar nota quando estiver nesse modo.
* As notas são preenchidas atributo por atributo.
* O sistema salva automaticamente as notas **apenas quando o usuário preencher todos os atributos de um item**.

  * Exemplo:
    Um item possui 5 atributos.
    O sistema só salva quando os 5 campos estiverem preenchidos.

---

## 7. Ordenação

O usuário pode ordenar os itens:

* Por qualquer atributo;
* Ou pelo **Overall**.

Regras:

* Apenas **um critério de ordenação por vez**;
* A ordenação pode ser crescente ou decrescente;
* A ordenação afeta apenas a visualização atual, não altera dados.

---

## 8. Visibilidade e Acesso

Cada ranking possui um tipo de visibilidade:

### Rankings Privados

* Só podem ser acessados por quem possui o **ID do ranking**;
* O ID é um identificador único no formato UUID v7;
* O acesso ocorre por URL direta, por exemplo:

  ```
  https://ranking.com/rankings/019bc9c7-b857-72c5-a491-7f9686c7989a
  ```
* Quem não possui essa URL não consegue encontrar o ranking.

### Rankings Públicos

* Podem ser visualizados por qualquer pessoa;
* Aparecem nos resultados de busca;
* Podem ser descobertos na página inicial.

---

## 9. Descoberta de Rankings

Na **página inicial**:

* São exibidos os **10 últimos rankings públicos criados**;
* O usuário pode buscar rankings:

  * Por parte do nome;
  * Por tags.

Exemplos de busca:

* “games” → encontra “Vídeo Games”, “Jogos de Tabuleiro”, etc.
* Tag “filmes” → encontra rankings marcados com essa tag.

---

## 10. Compartilhamento

Na página de cada ranking deve existir um botão de **Compartilhar**, com as opções:

* Copiar URL
* WhatsApp
* Instagram
* Facebook
* X (Twitter)

O objetivo é facilitar a divulgação rápida do ranking.

---

## 11. Identificação de Usuário (Sem Login)

O sistema **não possui autenticação tradicional** (login e senha).

Cada usuário recebe automaticamente um identificador único armazenado no navegador. Esse identificador é usado para:

* Saber quais rankings foram criados por aquele usuário;
* Saber quais notas aquele usuário deu em cada ranking.

O Product Owner está ciente de que esse método possui falhas (por exemplo, ao trocar de navegador ou limpar dados), mas isso é uma **decisão de negócio consciente**.
A prioridade é:

* Baixa fricção;
* Uso imediato;
* Nenhuma etapa de cadastro.

---

## 12. Meus Rankings

Sempre deve existir uma opção visível chamada **“Meus Rankings”**.

Essa área permite que o usuário:

* Veja todos os rankings que ele criou;
* Acesse rapidamente seus próprios rankings;
* Continue editando ou compartilhando-os.

---

## 13. Monetização

Todas as páginas do site devem conter espaços reservados para exibição de propagandas (ex.: Google Ads).

Esses espaços devem:

* Estar presentes em páginas principais (home, ranking, meus rankings);
* Não impedir o uso do sistema;
* Não interferir na leitura da tabela de ranking.

---

Este documento define **o comportamento do produto do ponto de vista do negócio**, servindo como referência única para qualquer IA ou pessoa envolvida na concepção da experiência do usuário.