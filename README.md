# Packer flasher postprocessor
Very simple postprocessor that dumps image onto given device with some sanity checks.

This plugin is used with [packer-builder-arm](https://github.com/mkaczanowski/packer-builder-arm) plugin to flash generated images to selected location.

# Configuration
```
"post-processors": [
 {
     "type": "flasher",
     "device": "/dev/sdX",
     "block_size": "4096",
     "interactive": true
 }
]   
```
