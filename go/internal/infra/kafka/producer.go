// Produzir mensagem no kafka
package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

type Producer struct {
	ConfigMap *ckafka.ConfigMap
}

func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{
		ConfigMap: configMap,
	}
}

// Plubicar mensagem pro kafka
// interface determina o formato da mensagem, nesse caso, a mensagem poderá ter qualquer formato,
// por isso definimos interface vazia {}
// key []byte = como uma string, mas em bytes
// key garante que a mesma mensagem caia na mesma partição
func (p *Producer) Publish(msg interface{}, key []byte, topic string) error {
	producer, err := ckafka.NewProducer(p.ConfigMap)

	// se houver algum erro, retornar o erro
	// por que as vezes nao se consegue produzir mas se consegue enviar a mensagem
	if err != nil {
		return err
	}

	// criar a mensagem usando o objeto do kafka &kafka.Message
	// TopicPartition = partição
	// nome do topico para onde a mensagem vai &topic
	// Partition: ckafka.PartitionAny , ou seja, kafka pode definir para qual partição ira a mensagem
	// msg.([]byte), força a mensagem a ser no formato byte
	message := &kafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		key:            key,
		Value:          msg.([]byte),
	}

	// produzir a mensagem:
	err = producer.Produce(message, nil)

	// se houver erro na produçao, retorna o erro:
	if err != nil {
		return err
	}

	// se nao houver nenhum erro, retorna erro em branco, por que é
	// obrigatorio enviar algum erro para a saída da func Publish:
	return nil

}
