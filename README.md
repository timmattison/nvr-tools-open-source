# nvr-tools-open-source

I'm releasing some closed source, commercial NVR tools eventually. In the meantime I'm putting some of the code here as
open-source so people can see what I'm doing and how I'm doing it.

I'm going to start with UniFi Protect and then move on to other NVRs if there's any interest.

The library code and tools here will be fairly bare bones. The commercial tools will have more features, GUI/TUI, etc.

## UniFi Protect

### Tools

#### list-cameras

The [list cameras](cmd/list-cameras) tool lists the cameras on a UniFi Protect NVR. It prints the camera MAC address,
type, name, and ID.

Install it with this command:

```bash
go install github.com/timmattison/nvr-tools-open-source/cmd/list-cameras@latest
```

#### list-license-plates

The [list license plates](cmd/list-license-plates) tool lists the license plate detections stored in a UniFi Protect
NVR. It prints the license plate number, detection start time, and detection end time.

Install it with this command:

```bash
go install github.com/timmattison/nvr-tools-open-source/cmd/list-license-plates@latest
```
