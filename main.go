package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Estrutura para o contrato JSON da mensagem
type Cliente struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type Metadados struct {
	Timestamp string `json:"timestamp"`
	Evento    string `json:"evento"`
}

type Mensagem struct {
	Cliente   Cliente   `json:"cliente"`
	Metadados Metadados `json:"metadados"`
}

func main() {
	// Carregar a configuração da AWS (ele utiliza automaticamente o arquivo de credenciais)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Criar um cliente SQS
	sqsClient := sqs.NewFromConfig(cfg)

	// URL da fila SQS
	queueURL := "https://sqs.sa-east-1.amazonaws.com/977099011826/customers"

	// Ler mensagens da fila
	for {
		// Definir os parâmetros de leitura da mensagem
		output, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 10,        // Lê até 10 mensagens por vez
			WaitTimeSeconds:     10,        // Long polling de 10 segundos
			VisibilityTimeout:   int32(30), // Define o tempo de visibilidade da mensagem em 30 segundos
		})
		if err != nil {
			log.Printf("erro ao ler mensagem da fila: %v", err)
			continue
		}

		// Processar cada mensagem
		for _, msg := range output.Messages {
			var mensagem Mensagem
			err := json.Unmarshal([]byte(*msg.Body), &mensagem)
			if err != nil {
				log.Printf("erro ao decodificar mensagem JSON: %v", err)
				continue
			}

			// Exibir a mensagem
			fmt.Printf("Cliente ID: %s, Nome: %s\n", mensagem.Cliente.ID, mensagem.Cliente.Nome)
			fmt.Printf("Evento: %s, Timestamp: %s\n", mensagem.Metadados.Evento, mensagem.Metadados.Timestamp)

			// Apagar a mensagem da fila após o processamento
			_, err = sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Printf("erro ao apagar a mensagem da fila: %v", err)
			}
		}

		// Intervalo entre as leituras
		time.Sleep(5 * time.Second)
	}
}
