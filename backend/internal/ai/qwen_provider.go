package ai

// QwenProvider is kept as a shim for backward compatibility.
// It wraps QwenClient and satisfies the legacy AIProvider interface.
// New code should use AIService which uses QwenClient directly.
// This file intentionally left minimal — all logic lives in qwen_client.go and service.go.
