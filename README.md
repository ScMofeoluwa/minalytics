# Minalytics

Minalytics is an **open-source**, **lightweight**, and **privacy-focused** web analytics tool designed to be simple alternative to platforms like Vercel Analytics. It provides essential insights into your website's traffic and user behavior without compromising user privacy. Minalytics is built to be self-hosted or used as a hosted service, giving you full control over your data.

----
## Why Minalytics?

Minalytics is designed with simplicity, privacy, and performance in mind. Here's what makes it stand out:

- **Open Source**: Fully transparent and customizable. You can self-host it or use the hosted service.
- **Privacy-Friendly**: No cookies, no unique identifiers, and no invasive tracking. Minalytics respects user privacy and complies with GDPR.
- **Unique Visits Tracking**: Uses a privacy-friendly approach (inspired by Vercel) to identify unique visits without cookies or persistent identifiers.
- **Custom Events**: Track custom events to monitor user interactions and behavior.
- **Hosted Service**: A user-friendly hosted platform with a dashboard for managing and analyzing your data.
- **Self-Hosted**: Run it on your own server for complete control over your data and infrastructure.

----
## Features
### Core Features

- **Unique Visits Tracking**: Identify unique visitors using a non-identifiable hash (no cookies or persistent identifiers).
- **Custom Events**: Track custom events to monitor specific user interactions on your website.
- **App-Based Tracking**: Create and manage multiple apps to track different websites or projects.
- **Geolocation**: Resolve user geolocation based on IP address.
- **Lightweight Integration**: Add Minalytics to your site with a simple script tag or integrate it into your backend.

### Analytics Insights

- **Page Views**: Track the number of views for each page.
- **Referrals**: Monitor where your traffic is coming from.
- **Devices**: Understand the types of devices your visitors are using.
- **Browsers**: Track browser usage statistics.
- **Operating Systems**: Monitor the operating systems used by your visitors.
- **Countries**: See where your visitors are located globally.
- **Visitors**: Get insights into unique and returning visitors.

### Hosted Service (Coming Soon)

- **Dashboard**: A user-friendly UI to manage your apps, view analytics, and export data.
- **Data Ownership**: Export your data at any time.

----
## Integration

To start tracking visits and events on your website, download the `minalytics.min.js` file from the Releases tab on GitHub. Replace `YOUR_TRACKING_ID` with the unique tracking ID for your app
```html
<script src="minalytics.min.js" tracking-id="YOUR_TRACKING_ID">
```

----
## Roadmap

- **Custom Events**: Enhance custom event tracking with more flexibility and metadata support.
- **Hosted Service**: Launch a hosted platform with a dashboard for managing apps and viewing analytics.
- **UI Dashboard**: Build a user-friendly interface for managing apps, viewing insights, and exporting data.

----
## Contributing

Minalytics is 100% open source, and we welcome contributions! Whether you're a developer, designer, or just someone with ideas, we'd love to have you on board.
