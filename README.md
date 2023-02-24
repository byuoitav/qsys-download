Microservice to download audio recordings from Q-Sys DSP's and upload them to Box

This microservice requires two files to be stored in ../ of the compiled code.  Namely, box_api_key.yourmom and box_folder_id.yourmom.  Their contents are a single lone of the api key and folder id from Box.

The http put should be formatted: https://serverIP/api/v1/coreIP/download/fileName
Content type needs to be "application/x-www-form-urlencoded"
two key-value pairs are required: filePath and room

filePath is the path to the file on the QSC core and should begin at and include the /Audio folder.

room is typically the prefix in the recorder, and is used to set a sub-folder in Box for organization of multiple recording "types".  We use the building-room naming convention