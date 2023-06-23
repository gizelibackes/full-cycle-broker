package kafka

// comando para baixa a dependencia
// & go mod tidy
// verificar no go.mod se a dependencia foi baixada corretamente
import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

type Consumer struct {
	ConfigMap *ckafka.ConfigMap
	Topics    []string
}

// Objeto para recepção do ConfigMap - com todas as configuracoes do kafka
func NewConsumer(configMap *ckafka.ConfigMap, topics []string) *Consumer {
	return &Consumer{
		ConfigMap: configMap,
		Topics:    topics,
	}
}

// Metodo de consumo das mensagens do kafka em multiplas treds
// as mensagens serão recebidas em um canal, e esse canal é enviado para o canal das ordens de compra e venda
// comunicaçao entre as treds através dos canais

func (c *Consumer) Consume(msgChan chan *ckafka.Message) error {
	consumer, err := ckafka.NewConsumer(c.ConfigMap)
	// caso haja erro para consumir a mensagem, aplicacao retorna um erro
	if err != nil {
		panic(err)
	}

	// caso haja problema para incrição nos tpics, aplicacao retorna erro
	err = consumer.SubscribeTopics(c.Topics, nil)
	if err != nil {
		panic(err)
	}

	// loop infinito para ler as mensagens que chegam no tipic, e enviar para o canal
	for {

		// se a mensagem tiver erro, fico lendo a mensagem:
		msg, err := consumer.ReadMessage(-1)

		// senao houver nenhum erro, a mensagem é enviada para o canal
		if err == nil {
			msgChan <- msg
		}
	}
}
