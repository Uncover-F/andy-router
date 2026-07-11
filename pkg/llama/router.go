package llama

import "github.com/charmbracelet/log"

func SelectModel(tps int) string {
	if tps < 300 {
		return ""
	} else if tps < 700 {
		log.Info("router: model selected", "model", "Mindcraft-CE/Andy-4.2-Micro-GGUF")
		return "Mindcraft-CE/Andy-4.2-Micro-GGUF"
	} else if tps < 1200 {
		log.Info("router: model selected", "model", "Mindcraft-CE/Andy-4.2-Air-GGUF")
		return "Mindcraft-CE/Andy-4.2-Air-GGUF"
	} else {
		log.Info("router: model selected", "model", "Mindcraft-CE/Andy-4.2-GGUF")
		return "Mindcraft-CE/Andy-4.2-GGUF"
	}
}
