# HandhelDB GUI  
A graphical interface designed with SDL for low-end devices.

## Compiling for TSP

To compile for TSP, refer to the instructions in the [CrossMix-OS toolchain README](https://github.com/cizia64/CrossMix-OS/blob/main/_assets/toolchain/README.md), and install:

## Installing GO for ARM64

Follow these steps to install GO for ARM64:

1. Update the package list:
   ```
   sudo apt-get update
   ```

2. Download the GO tarball:
   ```
   wget https://go.dev/dl/go1.21.0.linux-arm64.tar.gz
   ```

3. Extract the GO tarball:
   ```
   sudo tar -xvf go1.21.0.linux-arm64.tar.gz
   ```

4. Move GO to the `/usr/local` directory:
   ```
   sudo mv go /usr/local
   ```

5. Set the GOROOT environment variable:
   ```
   export GOROOT=/usr/local/go
   ```

6. Set the GOPATH environment variable:
   ```
   export GOPATH=$HOME/go
   ```

7. Add GO to your PATH:
   ```
   export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
   ```

8. Apply the changes to your profile:
   ```
   source ~/.profile
   ```

## Running the Application

To run the application, execute the following command:

```
go run main.go
```

## Configuration

The default configuration file (`config.json`) is structured as shown below. You can enable debugging logs, change the control type to keyboard, and adjust the screen resolution. To add collections, simply follow the existing pattern.

### Default `config.json`:

```json
{
    "logs": false,
    "control": {
        "type": "joystick"
    },
    "screen": {
        "width": 1280,
        "height": 720
    },
    "repositories": {
        "music": {
            "name": "Musics",
            "path": "/mnt/SDCARD/Roms/MUSIC",
            "extlist": [".mp3"],
            "collections": [
                {
                    "name": "geniesduclassique_vol3no01", // https://archive.org/details/geniesduclassique_vol3no01
                    "unzip": false
                },
                {
                    "name": "geniesduclassique_vol3no02", // https://archive.org/details/geniesduclassique_vol3no02
                    "unzip": false
                }
            ]
        }
    }
}
```

### Adding New Repositories:

For example, to add a new collection of Game Boy Advance (GBA) games, you can add a new repository as follows:

```json
"gba": {
    "name": "Game Boy Advanced",
    "path": "/mnt/SDCARD/Roms/GBA",
    "extlist": [".gba", ".zip"],
    "collections": [
        {
            "name": "mycollection_archiveorg_zip_format",
            "unzip": true
        },
        {
            "name": "geniesduclassique_vol3no02_gba_format",
            "unzip": false
        }
    ]
}
```
