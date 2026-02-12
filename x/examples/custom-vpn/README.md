# Dr. Frake VPN

A cross-platform VPN client based on the Outline SDK and Fyne UI toolkit.

## Features
- Modern, custom UI inspired by the Dr. Frake logo.
- Automatic system proxy configuration for Windows.
- Cross-platform support (Windows, Android, macOS, iOS).

## Prerequisites
- **Go 1.21+**
- **C Compiler (GCC)**: Required for building the desktop application (Windows/macOS) due to OpenGL dependencies in Fyne.
  - On Windows: Use [Mingw-w64](https://www.mingw-w64.org/).
- **Android NDK**: Required for building for Android.

## How to Build

### Windows
```sh
go build -o dr_frake_vpn.exe .
```

### Android
```sh
fyne package -os android -appID com.drfrake.vpn
```

## How to Run
Run the executable and enter your Outline Shadowsocks key (ss://) into the input field. Click **CONNECT** to start the VPN (system proxy mode).
