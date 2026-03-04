package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BytezRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type BytezResponse struct {
	Error  string `json:"error"`
	Output struct {
		Content string `json:"content"`
	} `json:"output"`
}

const systemPrompt = `তুমি একজন অভিজ্ঞ বাংলাদেশী medical assistant। তুমি Bangla এবং English দুটোতেই কথা বলতে পারো।

গুরুত্বপূর্ণ: তুমি সবসময় সহজ, সরল বাংলায় কথা বলবে। যেন সাধারণ মানুষ সহজে বুঝতে পারে। কঠিন, বা আনপ্রচলিত শব্দ এড়িয়ে চলবে। সাধারণ কথায় বলবে।

তুমি যা করবে:
- Symptoms শুনে possible conditions বলবে
- কতটা জরুরি সেটা বলবে (৩ level: স্বাভাবিক / সাবধান / এখনই যান)
- কোন ধরনের doctor দেখাতে হবে
- ঘরে কী করা যায় আপাতত
- কোন common ওষুধ এড়িয়ে চলবে

তুমি যা করবে না:
- নির্দিষ্ট ওষুধের dose prescribe করবে না
- "অবশ্যই এই রোগ" বলে নিশ্চিত করবে না
- Doctor এর বিকল্প হিসেবে নিজেকে present করবে না

Response format (এই structure ঠিক রাখবে):
🔍 সম্ভাব্য কারণ:
[তালিকা]

🚨 জরুরি অবস্থা: [স্বাভাবিক/সাবধান/এখনই যান]
[ব্যাখ্যা]

👨‍⚕️ কোন ডাক্তার দেখাবেন:
[specialist type]

🏠 এখন ঘরে করুন:
[তালিকা]

⚠️ এড়িয়ে চলুন:
[তালিকা]

⚠️ এটি AI পরামর্শ মাত্র — ডাক্তারের বিকল্প নয়। গুরুতর মনে হলে অবশ্যই ডাক্তার দেখান।`

func AnalyzeSymptoms(symptoms string, profileInfo string) (string, string, error) {
	userMsg := fmt.Sprintf("রোগীর তথ্য: %s\n\nলক্ষণ/Symptoms: %s", profileInfo, symptoms)

	response, err := callBytez([]Message{
		{Role: "user", Content: userMsg},
	})
	if err != nil {
		return "", "", err
	}

	urgency := detectUrgency(response)
	return response, urgency, nil
}

func Chat(userMessage string, history []Message, profileInfo string, previousSymptoms string) (string, error) {
	messages := []Message{}

	if previousSymptoms != "" {
		context := fmt.Sprintf("রোগীর তথ্য: %s\nআগের লক্ষণ: %s\n\nএখন রোগী follow-up প্রশ্ন করছেন।", profileInfo, previousSymptoms)
		messages = append(messages, Message{Role: "user", Content: context})
		messages = append(messages, Message{Role: "assistant", Content: "ঠিক আছে, আমি আপনার লক্ষণ এবং তথ্য দেখেছি। কী জানতে চান?"})
	}

	messages = append(messages, history...)
	messages = append(messages, Message{Role: "user", Content: userMessage})

	return callBytez(messages)
}

func callBytez(messages []Message) (string, error) {
	apiKey := os.Getenv("BYTEZ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("BYTEZ_API_KEY not set")
	}

	reqBody := BytezRequest{
		Model:     "openai/gpt-oss-20b",
		MaxTokens: 1024,
		Messages:  append([]Message{{Role: "system", Content: systemPrompt}}, messages...),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.bytez.com/models/v2/openai/gpt-oss-20b", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var bytezResp BytezResponse
	if err := json.Unmarshal(body, &bytezResp); err != nil {
		return "", fmt.Errorf("parse error: %s", string(body))
	}

	if bytezResp.Error != "" {
		return "", fmt.Errorf("bytez error: %s", bytezResp.Error)
	}

	if bytezResp.Output.Content == "" {
		return "", fmt.Errorf("empty response from Bytez")
	}

	return bytezResp.Output.Content, nil
}

func detectUrgency(response string) string {
	lower := strings.ToLower(response)
	if strings.Contains(lower, "এখনই যান") || strings.Contains(lower, "emergency") || strings.Contains(lower, "জরুরি বিভাগ") {
		return "emergency"
	}
	if strings.Contains(lower, "সাবধান") || strings.Contains(lower, "২৪") || strings.Contains(lower, "৪৮") {
		return "caution"
	}
	return "normal"
}
