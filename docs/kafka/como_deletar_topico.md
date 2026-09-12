## Entrar na linha de comando
docker compose exec kafka1 bash

## Deletar tópico
kafka-topics --bootstrap-server localhost:19092 --delete --topic payment_topic

## Recriar
kafka-topics --bootstrap-server localhost:19092 --create --topic payment_topic --partitions 1 --replication-factor 1

# Caso não delete

`Topic 'payment_topic' already exists` não é o AKHQ: o **ZooKeeper ainda tem o tópico**. O `--create --zookeeper` consulta `/brokers/topics/payment_topic`. Se esse znode existe, o create recusa — mesmo que a UI e o `--list` via broker não mostrem nada.

Isso costuma acontecer quando a deleção ficou a meio: logs da partição sumiram (AKHQ/broker “não veem” o tópico), mas o metadata no ZooKeeper ficou.

## 1. Confirme o estado

```bash
docker compose exec kafka1 kafka-topics --zookeeper zookeeper:2181 --list
docker compose exec kafka1 kafka-topics --zookeeper zookeeper:2181 --describe --topic payment_topic
```

Se `--list` mostrar `payment_topic` (às vezes com `marked for deletion`) ou o `--describe` funcionar, o create está certo em recusar.

No ZooKeeper:

```bash
docker compose exec zookeeper zkCli.sh

ls /brokers/topics
ls /admin/delete_topics
ls /config/topics
get /brokers/topics/payment_topic
```

## 2. Delete de novo e só então recrie

Pare os consumers (`antifraud`, `transaction`, `account`) se ainda estiverem no ar.

```bash
docker compose exec kafka1 kafka-topics --zookeeper zookeeper:2181 --delete --topic payment_topic
```

Espere o nome **sumir** de `--list` e de `/admin/delete_topics`. Aí sim o `--create`.

Não rode `--create` enquanto o `--delete` não tiver terminado.

## 3. Se o delete não limpar o ZooKeeper

Metadata órfão. Remova os znodes **só neste ambiente local**:

```bash
docker compose exec zookeeper zkCli.sh

rmr /admin/delete_topics/payment_topic
rmr /brokers/topics/payment_topic
rmr /config/topics/payment_topic
```

No broker, apague os dirs de log (volume `./data/kafka1/data`):

```bash
docker compose exec kafka1 ls /var/lib/kafka/data
# se existir algo como payment_topic-0:
docker compose exec kafka1 rm -rf /var/lib/kafka/data/payment_topic-0
```

Reinicie o Kafka e recrie o tópico.

## 4. Atalho local (o mais confiável)

Os dados estão em `./data/kafka1` e `./data/zookeeper`. Enquanto esses dirs existirem, o tópico “volta a existir” no metadata.

```bash
docker compose down
# apague ./data/kafka1 e ./data/zookeeper
docker compose up
```

Depois recrie o tópico (ou deixe o auto-create no **produce**, não no `Assign` dos consumers).

---

O Kafka agora também está instável: nos logs, `kafka1` não resolve `zookeeper:2181`. Você subiu com `docker compose start zookeeper kafka1 akhq` depois de um shutdown. Se o broker não estiver ligado ao ZooKeeper de forma estável, AKHQ fica vazio e o CLI ainda fala com o ZK antigo — daí a lista vazia + “already exists”.

Antes de insistir no create, confira se os dois containers estão healthy e na mesma rede:

```bash
docker compose ps
docker compose exec kafka1 ping -c 1 zookeeper
```

Resumo: **não precisa criar de novo até o znode desaparecer**. Ou `--delete` até `--list` limpar, ou limpe ZK + logs, ou wipe dos volumes `data/`.