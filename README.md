Microservice to download audio recordings from Q-Sys DSP's and upload them to Box

The Box upload portion microservice requires two files to be stored in /config in relation to the compiled code.  Namely, box_api_key.cfg and box_folder_id.cfg.  Their contents are a single line of the api key and folder id from Box.

If the cfg files do not exist or contain an empty string, the Box portion of this microservice will not run and the file will instead be sotred to the tmp_audio folder and not deleted.  The hope is this will be useful to others who do not need or want the upload to box feature.

Please note, the .cfg folders are to be stored in the /config/ directory "/config/box_api_key.cfg"

The http put should be formatted: https://serverIP:8013/api/v1/coreIP/download/fileName
Content type needs to be "application/x-www-form-urlencoded"
two key-value pairs are required: filePath and room

FilePath is the path to the file on the QSC core and should begin at and include the /Audio folder.

Room is typically the prefix in the recorder, and is used to set a sub-folder in Box for organization of multiple recording "types".  We use the building-room naming convention

Added support for QSC login.  A qsc_login.cfg file will need to be added to the /config folder

Examples of the config files can be found under the config_examples folder.  The files will need to be modified and the extensions changed to .cfg to work.