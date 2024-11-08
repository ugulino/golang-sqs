package models

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
