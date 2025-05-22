# NotSamsa Backend

Important: In development environments and due to restrictions of the user S3 minio-client regarding the PresignedURL generation with custom hostname
it's mandatory to set the following entry in your /etc/hosts

`127.0.0.1 s3`

Otherwise the download of preview images won't work.
