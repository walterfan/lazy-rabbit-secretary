# 🐰 Lazy Rabbit Secretary

Your AI-powered productivity companion to boost efficiency and never miss important tasks.

## 🚀 Overview

Lazy Rabbit Secretary is a comprehensive productivity management system that combines the power of AI with proven productivity methodologies like GTD (Getting Things Done) and Pomodoro Technique. Built with modern technologies, it provides an intuitive interface for managing tasks, reminders, and productivity analytics.

## ✨ Key Features

### 🍅 Pomodoro Timer
- 25-minute focused work periods followed by short breaks
- Enhanced focus and reduced distractions
- Better time awareness and productivity tracking

### 📋 GTD Flow Helper
- Complete Getting Things Done methodology implementation
- Capture, clarify, organize, reflect, and engage with tasks systematically
- Stress-free productivity with clear mind organization

### ✅ Smart Todo Lists
- Intelligent task management with priority levels and difficulty ratings
- Automatic scheduling based on your preferences
- Deadline tracking and progress visualization

### 🔔 Smart Reminders
- Never miss important tasks with intelligent notifications
- Customizable reminders based on urgency and daily patterns
- Multiple reminder types and smart scheduling

### 🤖 AI Agent Assistant
- Personal AI assistant for task automation and pattern analysis
- Smart suggestions and workflow optimization
- Learns from your behavior for personalized recommendations

### 📊 Productivity Insights
- Detailed analytics on your efficiency and task completion patterns
- Peak hours identification and productivity trends
- Performance tracking and habit formation support

## 🎯 How It Works

1. **Capture Tasks** - Quickly add tasks, ideas, and reminders to your inbox
2. **AI Organization** - Let AI categorize, prioritize, and schedule your tasks
3. **Focus & Execute** - Use Pomodoro sessions to focus and complete tasks efficiently
4. **Review & Improve** - Analyze your productivity and continuously optimize

## 🛠 Technology Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: Gin Gonic for HTTP routing
- **Database**: SQLite with GORM
- **AI Integration**: OpenAI API with function calling
- **Configuration**: Viper for YAML configuration
- **Logging**: Zap structured logging
- **Authentication**: JWT-based auth with OAuth support

### Frontend
- **Framework**: Vue 3 with Composition API
- **State Management**: Pinia
- **UI Framework**: Bootstrap 5 with custom styling
- **Build Tool**: Vite
- **Language**: TypeScript

### Additional Features
- **Email System**: Configurable email templates and notifications
- **Cron Jobs**: Automated task checking and reminder generation
- **Blog System**: AI-powered blog generation
- **Calendar Generation**: Automated calendar content creation
- **Secret Management**: Encrypted secret storage
- **BDD Testing**: Cucumber/Godog for behavior-driven development

## ⚡ Quick Start Guide

### Getting Started
1. Create your account and log in
2. Set up your profile and preferences
3. Add your first task or import existing ones
4. Configure notification settings
5. Start your first Pomodoro session!

### Pro Tips
- Use tags to categorize your tasks
- Set realistic deadlines and priorities
- Review your productivity analytics weekly
- Take regular breaks to maintain focus
- Use the AI assistant for task suggestions

## ❓ Frequently Asked Questions

### How does the Pomodoro timer work?
The Pomodoro timer follows the classic 25-minute work session followed by a 5-minute break. After 4 sessions, you get a longer 15-30 minute break. This helps maintain focus and prevents burnout.

### What is GTD methodology?
GTD (Getting Things Done) is a productivity methodology by David Allen. It involves capturing all tasks, clarifying what they mean, organizing them by context, reflecting on priorities, and engaging with them systematically.

### How does the AI assistant help with productivity?
The AI assistant analyzes your task patterns, suggests optimal scheduling, identifies productivity trends, and can automate routine tasks. It learns from your behavior to provide increasingly personalized recommendations.

## 🏗 Project Structure

```
lazy-rabbit-secretary/
├── backend/              # Go backend application
│   ├── cmd/             # Command-line interface
│   ├── config/          # Configuration files
│   ├── internal/        # Internal packages
│   │   ├── api/         # API handlers
│   │   ├── auth/        # Authentication logic
│   │   ├── jobs/        # Background job handlers
│   │   ├── models/      # Data models
│   │   ├── task/        # Task management
│   │   ├── reminder/    # Reminder system
│   │   └── post/        # Blog system
│   └── pkg/             # Shared packages
├── frontend/            # Vue.js frontend application
│   ├── src/
│   │   ├── components/  # Reusable Vue components
│   │   ├── views/       # Page components
│   │   ├── stores/      # Pinia state management
│   │   └── router/      # Vue Router configuration
├── deploy/              # Deployment configurations
└── doc/                 # Documentation
```

## 🚀 Development Status

- [x] Go backend with REST API
- [x] Vue.js frontend with modern UI
- [x] Task management system
- [x] Reminder system with email notifications
- [x] AI integration with OpenAI
- [x] User authentication and authorization
- [x] Blog system with AI content generation
- [x] Calendar generation
- [x] Secret management
- [x] Productivity analytics
- [ ] Mobile application
- [ ] Advanced AI features
- [ ] Integration with external calendar systems

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

Walter Fan - [GitHub](https://github.com/walterfan)

---

*Transform your productivity with AI-powered task management! 🚀*