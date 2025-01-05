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

## Trimui Smart Pro Installation

1. Download the latest release tagged with `trimui`.  
2. Extract the contents and locate the `HandheldDatabase` folder.  
3. Copy the `HandheldDatabase` folder to the `Apps` directory on your Trimui Smart Pro (TSP).  
4. Restart the device.  
5. Connect to Wi-Fi to complete the setup.

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

You will find a file called config.json in the config folder, open it, inside it there will be a list of repositories, one of them is music.
You can add as many as you want. Let's add a new one that points to this collection of DOS abandonwares: https://archive.org/details/Various_DOS_Abandonware_Ark

```json
{
   "logs":false,
   "control":{
      "type":"joystick"
   },
   "screen":{
      "width":1280,
      "height":720
   },
   "repositories":{
      "music":{
         "name":"Musics",
         "path":"/mnt/SDCARD/Roms/MUSIC",
         "extlist":[
            ".mp3"
         ],
         "collections":[
            {
               "name":"geniesduclassique_vol3no01",
               "unzip":false
            },
            {
               "name":"geniesduclassique_vol3no02",
               "unzip":false
            }
         ]
      },
      "dos":{
         "name":"DOS Games",
         "path":"/mnt/SDCARD/Roms/DOS",
         "extlist":[
            ".zip"
         ],
         "collections":[
            {
               "name":"Various_DOS_Abandonware_Ark" // you need just the collection name in url,
               "unzip":false
            }
         ]
      }
   }
}
```

To check if your JSON is valid, use the website: https://jsonformatter.curiousconcept.com/#

And just save it (remember that this config.json file must be in the tsp)
