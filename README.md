Microservice to download audio recordings from Q-Sys DSP's and upload them to Box

This microservice requires two files to be stored in ../ of the compiled code.  Namely, box_api_key.cfg and box_folder_id.cfg.  Their contents are a single lone of the api key and folder id from Box.

If the cfg files do not exist or contain an empty string, the Box portion of this microservice will not run and the file will instead be sotred to the tmp_audio folder and not deleted.  The hope is this will be useful to others who do not need or want the upload to box feature.

The http put should be formatted: https://serverIP:8013/api/v1/coreIP/download/fileName
Content type needs to be "application/x-www-form-urlencoded"
two key-value pairs are required: filePath and room

FilePath is the path to the file on the QSC core and should begin at and include the /Audio folder.

Room is typically the prefix in the recorder, and is used to set a sub-folder in Box for organization of multiple recording "types".  We use the building-room naming convention