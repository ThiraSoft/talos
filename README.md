# 🚀 Talos: Multi-Agent Generation AI 📖

## Warning : currently under heavy construction and should not be used in production.

## Overview

Talos is Go, AI-powered agentic framework that transforms tasks into an AI team workflow.
Leveraging the power of Google's Gemini AI, Talos orchestrates multiple specialized agents to collaboratively craft or do what you want them to do.

## 🌟 Features

- **Multi-Agent actions**: Collaborate between specialized AI agents.
- **Modular Design**: Easily extensible architecture for different tasks

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **AI Provider**: Google Gemini (Probably more coming soon)
- **Libraries**:
  - `google.golang.org/genai`
  - `github.com/google/uuid`

## 🚀 Quick Start

### Prerequisites

- Go 1.23.0+
- Google Gemini API Key

### Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/talos.git
cd talos
```

2. Set your Gemini API Key

```bash
export GEMINI_API_KEY=your_api_key_here
```

## 📖 Example Use Case

The current implementation focuses on generating an "Epic Fantasy" narrative, demonstrating the framework's capabilities.

## 🤖 Agents

Talos workflow exemple currently includes three core agents:

1. **AUTHOR**: Crafts the overall narrative and story structure
2. **WORLD_SHAPER**: Develops unique environments and world-building elements
3. **CHARACTER_SHAPER**: Creates complex, nuanced characters

## 🔧 Customization

Easily extend the framework by:

- Adding new agents
- Defining custom tasks
- Implementing additional tools

## 📄 License

This project is licensed under the MIT License. See the LICENSE file for details.
