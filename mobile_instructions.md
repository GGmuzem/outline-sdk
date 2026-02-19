# Dr. Frake Mobile App (Flutter + Go)

We will use **Flutter** for the UI and **Go Mobile** for the VPN logic.

## Architecture
1.  **Go Core**: `x/core` compiled to native libraries (`.aar`, `.framework`).
2.  **Native Bridge**: Small layer of Kotlin (Android) and Swift (iOS) that exposes Go functions.
3.  **Flutter UI**: Calls the Native Bridge via **MethodChannels**.

## Step 1: Generate Bindings
Run these commands in your project root:
```bash
cd x/core
gomobile bind -target=android -o ../../android_libs/DrFrakeCore.aar .
gomobile bind -target=ios -o ../../ios_libs/DrFrakeCore.framework .
```

## Step 2: Create Flutter Project
```bash
flutter create drfrake_mobile
```

## Step 3: Integrate Android
1.  Copy `DrFrakeCore.aar` to `drfrake_mobile/android/app/libs/`.
2.  Edit `android/app/build.gradle` to include the library.
3.  Edit `MainActivity.kt`:
    -   Initialize `DrFrakeCore`.
    -   Set up request handler for `MethodChannel("vpn_channel")`.
    -   Route calls like `login` and `connect` to the Go library.

## Step 4: Integrate iOS
1.  Drag `DrFrakeCore.framework` into Xcode project (`ios/Runner`).
2.  Edit `AppDelegate.swift`:
    -   Import `DrFrakeCore`.
    -   Set up `FlutterMethodChannel`.
    -   Route calls to Go library.

## Step 5: Flutter Logic (Dart)
```dart
static const platform = MethodChannel('vpn_channel');

Future<void> login(String email, String pass) async {
  await platform.invokeMethod('login', {'email': email, 'pass': pass});
}

Future<void> connect(String config) async {
  await platform.invokeMethod('connect', {'config': config});
}
```
