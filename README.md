# 🚀 ads-zero

**ads-zero** is an alpha-stage alerting system built for modern marketers. It fetches campaign data from channels like **Facebook**, **Google**, **TikTok**, and **Taboola**, applies customizable rules, and sends notifications via **email**, **Telegram**, or **Slack** when conditions are met.

> **Note:** Currently, only the Facebook integration is implemented, with ClickHouse used for data storage. Postgres support is planned along with additional data sources. The Kafka engine for scalable deployment is fully implemented.

## ✨ Features

| Component                     | Status                          | Details                                                                                   |
| ----------------------------- | ------------------------------- | ----------------------------------------------------------------------------------------- |
| **Data Fetching**             | Facebook ✅ <br> Google ⏳ <br> TikTok ⏳ <br> Taboola ⏳ | Fetches marketing campaign data from multiple sources.                                    |
| **Rule Engine**               | Customizable ✅                 | Execute one or more rules against the fetched data to detect defined conditions.          |
| **Notification**              | Email, Telegram, Slack ✅        | Send alerts to users when specific conditions are met.                                    |
| **Data Storage**              | ClickHouse ✅ <br> Postgres ⏳      | Save the fetched data for further analysis.                                               |
| **Deployment**                | Single Component ✅ <br> Kafka ✅  | Use as a standalone component with an internal ticker or deploy at scale using Kafka.      |
| **Configuration**             | Database Configurable ✅        | Configure alerts via database table records. Planned: JSON interface for simpler setups.   |

## 📊 Project Status

- **Alpha Stage:** Minimal implementation for Facebook data integration.
- **Storage:** ClickHouse is the current working solution. Postgres support is planned.
- **Deployment:** Supports both single component mode (internal ticker) and scalable deployment using Kafka.
- **Configuration:** Alerts are currently configured through database table records. A JSON-based configuration interface is planned for users who do not need long-term data storage.

## 🚀 Getting Started

1. **Clone the Repository**
   ```sh
   git clone <repository_url>
   cd ads-zero-1
   ```
2. Install any necessary dependencies.
3. Configure your notification channels (email, Telegram, or Slack).
4. Start testing and contributing to the development of additional data sources and rule implementations.

## Contribution

Contributions are welcome! Please feel free to fork the repository and submit pull requests with improvements or new features.

## License

This project is licensed under the terms of the Apache License 2.0.

## Contact

For any questions or collaboration ideas, please open an issue in the repository.