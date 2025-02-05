---
title: "Receipt Splitter"
date: "2025-02-04"
description: "AI-powered receipt splitting"
github: "https://github.com/lewislewin/receipt-splitter"
---

# Receipt Splitter

## What is it?

Splitting a bill at a restaurant or pub is a pain. **Who owes what? How do we fairly divide shared items?**  
The **Receipt Splitter App** solves this with a **fully automated, AI-powered system**.

## How It Works

1. **Upload a Receipt** 📸 – The person who paid uploads a photo of the receipt.
2. **AI Parses the Receipt** 🤖 – The app uses **Google OCR** to extract text from the image and then sends it to **OpenAI** to convert it into a structured JSON object, identifying items, prices, tax, and deductions.
3. **Share the Link** 🔗 – The uploader sends a link to everyone at the table, providing access to a digital version of the receipt.
4. **Select Your Items** 👥 – Each person selects what they ordered.
5. **Monzo Links Generated** 💳 – Instantly generate a Monzo payment link for seamless reimbursement.

## Why It's Useful

- **No more manual calculations** – AI handles all the math.  
- **Fair & Transparent** – Everyone picks their items before paying.  
- **Works with any receipt** – OCR ensures accuracy.  
- **Monzo Integration** – Send payments with one tap.

## The Tech Behind It

- **Frontend:** Built with **SvelteKit** for a smooth user experience.
- **Backend:** Initially written in **Django**, but Go provided better **performance and simplicity**.
- **AI Parsing:** Utilizes **Google OCR** to extract text from receipt images, which is then processed by **OpenAI** to convert it into a structured JSON object.
- **Payments:** Generates **Monzo.me links** dynamically.

*Note: The project's source code is available on GitHub: [lewislewin/receipt-splitter](https://github.com/lewislewin/receipt-splitter) and [lewislewin/receipt-splitter-backend](https://github.com/lewislewin/receipt-splitter-backend).*
