# s3server

Run a webserver serving files from AWS S3.

## Usage

### Serve assets from the root of S3bucket

```
s3server yourbucketname
```

Note that `yourbucketname` is just the name of the bucket, e.g. `dev-s3server`.

### Serve assets under a folder of S3bucket

```
s3server yourbucketname foldername
```
