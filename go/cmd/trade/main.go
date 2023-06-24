package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gizelibackes/Full-cycle-HomeBroker/go/internal/infra/kafka"
	"github.com/gizelibackes/Full-cycle-HomeBroker/go/internal/market/dto"
	"github.com/gizelibackes/Full-cycle-HomeBroker/go/internal/market/entity"
	"github.com/gizelibackes/Full-cycle-HomeBroker/go/internal/market/transformer"
)

// definir entrypoint da aplicação, ou seja,
// definir canal de entrada pra receber os dados da order
// defnir canal de saida para retonar os dados da order
func main() {
	// por onde recebemos os dados via kafka
	ordersIn := make(chan *entity.Order)
	// por onde enviamos os dados via kafka
	ordersOut := make(chan *entity.Order)
	// wait Group
	wg := &sync.WaitGroup{}
	// defer para executar o wg.Wait() so no final do codigo
	defer wg.Wait()
	// canal para receber informaçoes do kafka
	// ckafka = apelido definido para o pacote oficial do kafka
	kafkaMsgChan := make(chan *ckafka.Message)
	// criar coneccao com o kafka atraves do docker
	configMap := &ckafka.ConfigMap{
		// alterar etc/host no Mac:
		// $ sudo nano /etc/hosts
		// 127.0.0.1       host.docker.internal
		"bootstrap.services": "host.docker.internal:9094",
		"group.id":           "myGroup",
		"auto.offset.reset":  "earliest",
	}

	//criar  producer e consumer
	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	// criar uma tread de solicitacao de consumo das mensagens do kafka
	go kafka.Consume(kafkaMsgChan) // T2

	// criar o book passando canal de entrada e canal de saida
	// recebe do canal do kafka, joga no input, processa joga no output e depois publica no kafka
	book := entity.NewBook(ordersIn, ordersOut, wg)

	// criar tread para receber os dados das mensagens
	go book.Trade() // T3

	// criar uma tread de look infinito para receber os dados do kafka e publicar no input
	go func() {
		// le a mensagem e envia para o canal kafkaMsgChan
		for msg := range kafkaMsgChan {
			// para cada transacao, reinicar o wg
			wg.Add(1)
			fmt.Println(string(msg.Value))
			// transforma a msg de jason (json.Unmarshal) para o objeto
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err != nil {
				panic(err)
			}
			// transforma o dado para o formato de objeto
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	// pegar dados do canal que o book processou, transformar em dado cru,
	// e publicar novamente no kafka
	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		// transforma de objeto para json json.MarshalIndent
		outputJson, err := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}

		// enviar o dado para o kafka
		// mensagem = outputJson
		// key = []byte("orders")
		producer.Publish(outputJson, []byte("orders"), "output")
	}

}
