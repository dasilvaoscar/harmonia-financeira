# 4 - Transactional Outbox Pattern

**Título da Decisão: Testes de estresse**

## 1. Contexto
Como medida de resiliência entre a publicação de eventos, essa ADR propõe a utilização do padrão
Transactional Outbox para manter a consistência na emissão dos eventos e alteração/criação de entidades.

## 2. Decisão
Cada serviço deve conter uma tabela `outbox` que será consumida via worker. O worker é executado a cada 10 segundos via cron.
Uma evolução da checagem dessa tabela é também fazer o consumo via CDC e o worker atuando como uma dupla validação,
para circunstâncias onde o CDC possa ter falhado.

id    |payload|entity|domainId
------|------|------|------|
uuid-v7|{ ...payload }|"ENUM"|uuid-v7

## **3. Consequências**  

**Positivas:**  
- Atomicidade local: Salva o dado principal e o evento na tabela outbox na mesma transação do banco.
- Consistência forte: Evita que o sistema perca eventos se o broker de mensagens cair.
- Fuga do 2PC: Elimina a necessidade de protocolos de confirmação em duas etapas (two-phase commit), que são lentos e complexos.
- Confiabilidade: Permite reprocessar mensagens pendentes de forma simples em caso de falha temporária.  

**Negativas:**  
- Possível necessidade de ajustes no ambiente de testes para suportar a execução dos testes de carga de forma eficiente.  
- Necessidade de lidar com um alto volume de escrita em duas tabelas diferentes de forma atômica.
- Checagem constante dos workers para envio das mensagens ao broker obriga a criação de um novo serviço desacoplado.
- O volume de dados da tabela outbox pode crescer rápidamente, obrigando a termos algum mecanismo de expurgo de dados.