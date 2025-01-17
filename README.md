# nativenotify

`nativenotify` provides a cross platform notification API for Go.

It aims to provide a convenient to use API covering the basic use cases.
Not all features for each platform is supported. 

Notably, text inputs are not supported on Linux.

If you need more power on a given platform, use the platform specific modules directly.

## Caveats

MacOS programs need to be bundled and codesigned for notifications to be allowed. 

Linux programs will get access to a better platform API if sandboxed, such as in a Flatpak.

Windows programs might run into issues depending on the version of Windows being used as not all versions support this notification API.

## Architecture

The package offers two high-level functions:

- `Setup`
- `Push`

All notification APIs need some initialization at the platform layer. The initialization requirements are completely different.
The configuration is therefore a "fat union" structure, which contains all the fields necessary for each platform.

`Setup` should be called exactly once before pushing notifications.

Notifications are described in a common form via the exported `Notification`, `TextAction` and `ButtonAction` types.
These provide lowest-common-denominator features.

The `ID` fields name the parent notification and the associated actions. 
The callback of the notification will always receive the parent notification ID, and a map of user data which will contain
the data for the action that was activated. This map will also include the parent notification payload.

The map needs to be inspected to understand what part of the notification was activated. 
