package kafka

func (ok *OrchestratorKafkaImpl) GetMaxRetries() int {
	return ok.MaxRetries
}
