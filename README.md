# Gostockly

Managing inventory across multiple Shopify stores can be challenging. To address this, I developed Gostockly, which ensures real-time synchronization of stock levels across all connected stores.

## Key Features

- **Prevents Overselling:** By keeping inventory levels consistent, the tool helps avoid situations where products are sold beyond available stock.
- **Automatic Updates:** The system detects and applies stock changes instantly, ensuring all stores reflect accurate inventory without manual intervention.
- **User-Friendly Dashboard:** A centralized interface allows for easy monitoring and management of stock levels across all stores.

## Technical Implementation

The application is built using Go, chosen for its performance and concurrency capabilities. It leverages Shopify's API to manage inventory data. The project structure includes:

- **cmd/**: Contains the main application entry point.
- **config/**: Manages configuration settings.
- **docker/**: Includes Docker-related files for containerization.
- **internal/**: Houses internal application logic.
- **pkg/**: Contains reusable packages.

## Frontend Integration

To provide a seamless user experience, I developed a frontend application using SvelteKit, a modern framework for building web applications. This frontend communicates with the Go backend to display real-time inventory data and synchronization status. The repository for the frontend is available here: [lewislewin/gostockly-frontend](https://github.com/lewislewin/gostockly-frontend).

## Importance and Backstory

In the rapidly evolving e-commerce landscape, businesses often operate multiple online stores to reach diverse markets or manage different product lines. However, this multi-store approach introduces the challenge of maintaining consistent inventory levels across all platforms. Discrepancies can lead to overselling, stockouts, and dissatisfied customers.

Recognizing this pain point, I developed Gostockly to provide a reliable solution for real-time inventory synchronization. By ensuring that all connected stores reflect accurate stock levels, businesses can prevent overselling, optimize inventory management, and enhance customer satisfaction.

## Getting Started

For those interested in implementing this solution or seeking further information, please reach out via email at [lewis@lewislewin.com](mailto:lewis@lewislewin.com).

*Note: The project's source code is available on GitHub: [lewislewin/gostockly](https://github.com/lewislewin/gostockly) and  [lewislewin/gostockly-frontend](https://github.com/lewislewin/gostockly-frontend).*
