# 🏥 AI ডাক্তার সহায়ক

বাংলাদেশের জন্য AI-চালিত প্রাথমিক স্বাস্থ্য পরামর্শ সিস্টেম।
(Bytez API) + Go + SQLite + Vanilla HTML/CSS/JS

---

## 📁 Project Structure

```
ai-doctor-bd/
├── backend/
│   ├── main.go
│   ├── go.mod
│   ├── .env.example
│   ├── db/
│   │   ├── db.go
│   │   └── schema.sql
│   ├── handlers/
│   │   ├── profile.go
│   │   └── consultation.go
│   ├── models/
│   │   └── models.go
│   └── ai/
│       └── bytez.go
└── frontend/
    ├── index.html      ← প্রধান পরামর্শ পাতা
    ├── history.html    ← ইতিহাস পাতা
    ├── profiles.html   ← প্রোফাইল ব্যবস্থাপনা
    └── css/
        └── style.css
```

---

## 🚀 চালু করার নিয়ম

### ১. Backend Setup

```bash
cd backend

# .env ফাইল তৈরি করুন
cp .env.example .env
# .env ফাইলে আপনার Bytez API key দিন

# Dependencies ইনস্টল করুন
go mod tidy

# Backend চালু করুন
go run main.go
```

Backend চালু হলে দেখাবে:
```
✅ Database initialized
🚀 AI Doctor BD backend running on http://localhost:8080
```

### ২. Frontend চালু করুন

যেকোনো browser এ `frontend/index.html` ফাইলটি খুলুন।

অথবা simple HTTP server দিয়ে:
```bash
cd frontend
python3 -m http.server 3000
# তারপর http://localhost:3000 খুলুন
```

---

## ⚙️ Environment Variables

```env
BYTEZ_API_KEY=your_bytez_api_key_here
```

Bytez API key পেতে: https://bytez.com

---

## 🔌 API Endpoints

| Method | Route | কাজ |
|--------|-------|-----|
| GET | `/api/profiles` | সব প্রোফাইল দেখুন |
| POST | `/api/profiles` | নতুন প্রোফাইল যোগ |
| DELETE | `/api/profiles/{id}` | প্রোফাইল মুছুন |
| POST | `/api/analyze` | AI লক্ষণ বিশ্লেষণ |
| POST | `/api/chat` | Follow-up chat |
| GET | `/api/history/{profileId}` | পরামর্শের ইতিহাস |
| GET | `/api/messages/{consultationId}` | Chat ইতিহাস |

---

## ✨ Features

- 🇧🇩 বাংলা ও ইংরেজি উভয়ে কাজ করে
- 👥 পরিবারের একাধিক সদস্যের প্রোফাইল
- 🤖 Qwen 3 AI দিয়ে লক্ষণ বিশ্লেষণ
- 🚨 ৩ স্তরের জরুরি অবস্থা নির্দেশক
- 💬 Follow-up chat সুবিধা
- 📋 সম্পূর্ণ ইতিহাস সংরক্ষণ
- ⚠️ সর্বদা Disclaimer সহ

---

## ⚠️ Disclaimer

এই tool টি **প্রাথমিক পরামর্শের** জন্য মাত্র।
গুরুতর অসুস্থতায় অবশ্যই একজন qualified ডাক্তারের পরামর্শ নিন।
