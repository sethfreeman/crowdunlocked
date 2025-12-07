# Crowd Unlocked Mobile

Flutter mobile application for Crowd Unlocked artist management platform.

## Getting Started

### Prerequisites

- Flutter SDK 3.2.0 or higher
- Xcode (for iOS development)
- Android Studio (for Android development)

### Installation

```bash
flutter pub get
```

### Running

```bash
# Run on connected device
flutter run

# Run on specific device
flutter devices
flutter run -d <device-id>
```

### Building

```bash
# Android
flutter build apk
flutter build appbundle

# iOS
flutter build ios
```

## Architecture

- **State Management**: Provider
- **HTTP Client**: http package
- **Local Storage**: shared_preferences

## Features

- Artist bookings management
- Music releases tracking
- Social media monitoring
- Revenue analytics
